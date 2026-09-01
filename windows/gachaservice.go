package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/scan252/aMC/windows/internal/gacha"
	"github.com/scan252/aMC/windows/internal/wulog"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// GachaService 抽卡域对前端的唯一入口。
type GachaService struct{}

// AccountSummary 账号列表项。
type AccountSummary struct {
	UID       string  `json:"uid"`
	SvrArea   string  `json:"svrArea"`
	FetchedAt string  `json:"fetchedAt"`
	Total     int     `json:"total"`
	Count5    int     `json:"count5"`
	AvgPity   float64 `json:"avgPity"`
}

// AccountDetail 打开的账号详情。
type AccountDetail struct {
	UID       string             `json:"uid"`
	SvrArea   string             `json:"svrArea"`
	FetchedAt string             `json:"fetchedAt"`
	Stats     gacha.OverallStats `json:"stats"`
	Dist      map[string][]int   `json:"dist"`
	Recent5   []RecentItem       `json:"recent5"`
}

// RecentItem 最近 5★ 列表项；Gap 为与上一个 5★ 的间隔（-1 表示窗口内无前序 5★）。
type RecentItem struct {
	gacha.Record
	Pool string `json:"pool"`
	Gap  int    `json:"gap"`
}

// LogCandidate 日志候选（供前端展示发现结果）。
type LogCandidate = wulog.LogInfo

// fetchEvent 抓取进度事件负载。
type fetchEvent struct {
	Index int    `json:"index"`
	Total int    `json:"total"`
	Pool  string `json:"pool"`
	Err   string `json:"err,omitempty"`
}

func emitProgress(index, total int, pool string, err error) {
	app := application.Get()
	if app == nil {
		return
	}
	e := fetchEvent{Index: index, Total: total, Pool: pool}
	if err != nil {
		e.Err = err.Error()
	}
	app.Event.Emit("gacha:progress", e)
}

// DiscoverLogs 自动发现本机鸣潮客户端日志。
func (s *GachaService) DiscoverLogs() ([]LogCandidate, error) {
	return wulog.Discover()
}

// ListAccounts 列出本地全部账号摘要。
func (s *GachaService) ListAccounts() ([]AccountSummary, error) {
	dir, err := gacha.DefaultDataDir()
	if err != nil {
		return nil, err
	}
	ids := gacha.ListPlayerIDs(dir)
	out := make([]AccountSummary, 0, len(ids))
	for _, id := range ids {
		data, err := gacha.LoadGachaData(gacha.PlayerDataPath(dir, id))
		if err != nil || data == nil {
			continue
		}
		stats := gacha.ComputeStats(data)
		out = append(out, AccountSummary{
			UID: id, SvrArea: data.SvrArea, FetchedAt: data.FetchedAt,
			Total: stats.Total, Count5: stats.Count5, AvgPity: stats.AvgPity,
		})
	}
	return out, nil
}

// OpenAccount 打开指定账号的完整数据与统计。
func (s *GachaService) OpenAccount(uid string) (*AccountDetail, error) {
	dir, err := gacha.DefaultDataDir()
	if err != nil {
		return nil, err
	}
	data, err := gacha.LoadGachaData(gacha.PlayerDataPath(dir, uid))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("账号 %s 无本地数据", uid)
	}
	return buildDetail(data), nil
}

// Refresh 完整抓取流程：发现日志 → 提取 URL → 拉取全部卡池 → 合并 → 存储。
// logPath 为空时走自动发现；playerID 为空时取日志中最近一条。
func (s *GachaService) Refresh(logPath, playerID string) (*AccountDetail, error) {
	if logPath == "" {
		logs, err := wulog.Discover()
		if err != nil || len(logs) == 0 {
			return nil, fmt.Errorf("未找到鸣潮客户端日志，请先在游戏中打开「唤取 → 唤取记录」页面，或在设置中手动指定日志路径")
		}
		logPath = logs[0].Path
	}

	_, rawURL, err := wulog.ParseLog(logPath, playerID)
	if err != nil {
		return nil, err
	}
	creds := gacha.ParseCredentials(rawURL)
	if creds.PlayerID == "" {
		return nil, fmt.Errorf("唤取 URL 中缺少 player_id")
	}

	pools := gacha.FetchAllPools(creds, func(i, total int, name string, err error) {
		emitProgress(i, total, name, err)
	})

	dir, err := gacha.DefaultDataDir()
	if err != nil {
		return nil, err
	}
	path := gacha.PlayerDataPath(dir, creds.PlayerID)
	existing, err := gacha.LoadGachaData(path)
	if err != nil {
		existing = nil
	}

	data := &gacha.GachaData{
		PlayerID:  creds.PlayerID,
		SvrArea:   creds.SvrArea,
		FetchedAt: time.Now().Format("2006-01-02 15:04:05"),
		Pools:     make(map[string][]gacha.Record, len(pools)),
	}
	for key, incoming := range pools {
		var old []gacha.Record
		if existing != nil {
			old = existing.Pools[key]
		}
		data.Pools[key] = gacha.MergePoolRecords(old, incoming)
	}
	if err := gacha.SaveGachaData(path, data); err != nil {
		return nil, err
	}
	return buildDetail(data), nil
}

// ExportAccount 把账号数据导出为带时间戳的 JSON 文件，返回目标路径。
func (s *GachaService) ExportAccount(uid, destDir string) (string, error) {
	dir, err := gacha.DefaultDataDir()
	if err != nil {
		return "", err
	}
	src := gacha.PlayerDataPath(dir, uid)
	raw, err := readFileBytes(src)
	if err != nil {
		return "", err
	}
	if destDir == "" {
		destDir = dir
	}
	dest := fmt.Sprintf("%s\\amc_%s_%s.json", destDir, uid, time.Now().Format("20060102_150405"))
	if err := writeFileBytes(dest, raw); err != nil {
		return "", err
	}
	return dest, nil
}

// ImportAccount 从 JSON 文件导入账号数据（兼容 Mac 版 aMC 导出格式），返回 UID。
func (s *GachaService) ImportAccount(path string) (string, error) {
	data, err := gacha.LoadGachaData(path)
	if err != nil {
		return "", err
	}
	if data == nil || data.PlayerID == "" {
		return "", fmt.Errorf("导入文件缺少 player_id，且不是有效的 aMC 数据文件")
	}
	if data.Pools == nil {
		data.Pools = map[string][]gacha.Record{}
	}
	dir, err := gacha.DefaultDataDir()
	if err != nil {
		return "", err
	}
	if err := gacha.SaveGachaData(gacha.PlayerDataPath(dir, data.PlayerID), data); err != nil {
		return "", err
	}
	return data.PlayerID, nil
}

func buildDetail(data *gacha.GachaData) *AccountDetail {
	d := &AccountDetail{
		UID:       data.PlayerID,
		SvrArea:   data.SvrArea,
		FetchedAt: data.FetchedAt,
		Stats:     gacha.ComputeStats(data),
		Dist:      gacha.PityDist(data),
		Recent5:   recentFiveStars(data, 20),
	}
	return d
}

// recentFiveStars 跨卡池取最近 n 条 5★ 记录，并计算与上一个 5★ 的间隔。
func recentFiveStars(data *gacha.GachaData, n int) []RecentItem {
	var items []RecentItem
	for key, records := range data.Pools {
		asc := make([]gacha.Record, len(records))
		copy(asc, records)
		sort.SliceStable(asc, func(i, j int) bool { return asc[i].Time < asc[j].Time })

		sinceLast := 0
		found5 := false
		for _, r := range asc {
			if r.QualityLevel != 5 {
				sinceLast++
				continue
			}
			gap := -1
			if found5 {
				gap = sinceLast + 1
			}
			items = append(items, RecentItem{Record: r, Pool: key, Gap: gap})
			found5 = true
			sinceLast = 0
		}
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].Time > items[j].Time })
	if len(items) > n {
		items = items[:n]
	}
	return items
}

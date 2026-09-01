package gacha

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseCredentialsCN(t *testing.T) {
	u := "https://aki-gm-resources.aki-game.com/aki/gacha/index.html#/record?record_id=R1&player_id=100252731&svr_id=01&lang=zh-Hans&resources_id=P1"
	c := ParseCredentials(u)
	if c.RecordID != "R1" || c.PlayerID != "100252731" || c.ServerID != "01" || c.CardPoolID != "P1" {
		t.Fatalf("国服凭证解析错误: %+v", c)
	}
	if c.IsGlobal {
		t.Fatal("com 域名不应判定为国际服")
	}
	if c.APIBase() != APIBaseCN {
		t.Fatalf("国服 API base 错误: %s", c.APIBase())
	}
}

func TestParseCredentialsGlobal(t *testing.T) {
	u := "https://aki-gm-resources-oversea.aki-game.net/aki/gacha/index.html#/record?record_id=R2&player_id=999888777&svr_id=usa&gacha_id=G2"
	c := ParseCredentials(u)
	if !c.IsGlobal {
		t.Fatal("oversea+.net 应判定为国际服")
	}
	if c.CardPoolID != "G2" {
		t.Fatalf("resources_id 缺失时应回退 gacha_id: %+v", c)
	}
	if c.SvrArea != "global" {
		t.Fatalf("svr_area 默认应为 global: %+v", c)
	}
	if c.APIBase() != APIBaseGlobal {
		t.Fatalf("国际服 API base 错误: %s", c.APIBase())
	}
}

func makeRec(name, time string, q int) Record { //nolint:revive
	return Record{CardPoolType: "1", ResourceID: 1, QualityLevel: q, ResourceType: "角色", Name: name, Count: 1, Time: time}
}

func TestMergePoolRecords(t *testing.T) {
	existing := []Record{
		makeRec("今汐", "2026-05-01 10:00:00", 5),
		makeRec("维里奈", "2026-04-01 10:00:00", 5),
	}
	incoming := []Record{
		makeRec("今汐", "2026-05-01 10:00:00", 5), // 与已有重复，应去重
		makeRec("椿", "2026-08-01 10:00:00", 5),
	}
	merged := MergePoolRecords(existing, incoming)
	if len(merged) != 3 {
		t.Fatalf("合并后应 3 条（去重+保留过期），got %d: %+v", len(merged), merged)
	}
	if merged[0].Name != "椿" || merged[2].Name != "维里奈" {
		t.Fatalf("应按时间降序: %+v", merged)
	}
}

func TestMergeKeepsExpiredHistory(t *testing.T) {
	expired := makeRec("安可", "2025-01-01 10:00:00", 5)
	merged := MergePoolRecords([]Record{expired}, []Record{makeRec("长离", "2026-08-01 10:00:00", 5)})
	if len(merged) != 2 {
		t.Fatalf("过期历史应保留，got %d", len(merged))
	}
}

func TestStoreRoundTripAndBackup(t *testing.T) {
	dir := t.TempDir()
	path := PlayerDataPath(dir, "100252731")

	data := &GachaData{
		PlayerID:  "100252731",
		SvrArea:   "cn",
		FetchedAt: "2026-09-01 12:00:00",
		Pools:     map[string][]Record{"1": {makeRec("今汐", "2026-05-01 10:00:00", 5)}},
	}
	if err := SaveGachaData(path, data); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadGachaData(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded == nil || loaded.PlayerID != "100252731" || len(loaded.Pools["1"]) != 1 {
		t.Fatalf("回读不符: %+v", loaded)
	}

	// 修改后再次保存应产生备份
	loaded.Pools["1"][0].Name = "改"
	if err := SaveGachaData(path, loaded); err != nil {
		t.Fatal(err)
	}
	backupDir := filepath.Join(dir, "100252731", "backup")
	entries, _ := os.ReadDir(backupDir)
	if len(entries) != 1 {
		t.Fatalf("应产生 1 份备份，got %d", len(entries))
	}

	ids := ListPlayerIDs(dir)
	if len(ids) != 1 || ids[0] != "100252731" {
		t.Fatalf("ListPlayerIDs 不符: %v", ids)
	}
}

func TestComputeStats(t *testing.T) {
	pool1 := []Record{
		makeRec("三星", "2026-08-01 10:00:00", 3),
		makeRec("四星", "2026-08-02 10:00:00", 4),
		makeRec("鉴心", "2026-08-03 10:00:00", 5), // 第 3 抽出金
		makeRec("三星", "2026-08-04 10:00:00", 3),
		makeRec("赞妮", "2026-08-05 10:00:00", 5), // 第 2 抽出金（间隔 2）
		makeRec("三星", "2026-08-06 10:00:00", 3), // 已垫 1
	}
	data := &GachaData{PlayerID: "u", Pools: map[string][]Record{"1": pool1}}
	stats := ComputeStats(data)

	if stats.Total != 6 || stats.Count5 != 2 || stats.Count4 != 1 {
		t.Fatalf("总览统计错误: %+v", stats)
	}
	if stats.AvgPity != 3.0 { // 总抽数 6 / 5★数 2
		t.Fatalf("平均出金应为 3.0，got %v", stats.AvgPity)
	}
	if stats.LuckIndex != ExpectedAvgPity/3.0 {
		t.Fatalf("欧非指数错误: %v", stats.LuckIndex)
	}

	p := stats.Pools[0]
	if p.Pity != 1 || p.PityIsFloor {
		t.Fatalf("当前垫抽数应为 1: %+v", p)
	}
	if p.Last5Name != "赞妮" {
		t.Fatalf("最近 5★ 应为赞妮: %+v", p)
	}

	dist := PityDist(data)["1"]
	if dist[0] != 1 { // 唯一相邻间隔样本：3★,5★ 间为 2 抽
		t.Fatalf("分布桶错误: %v", dist)
	}
}

func TestComputeStatsNoFiveStar(t *testing.T) {
	data := &GachaData{Pools: map[string][]Record{
		"1": {makeRec("三星", "2026-08-01 10:00:00", 3), makeRec("三星", "2026-08-02 10:00:00", 3)},
	}}
	stats := ComputeStats(data)
	if stats.AvgPity != 0 || stats.LuckIndex != 0 {
		t.Fatalf("无 5★ 样本时不应给出均值/指数: %+v", stats)
	}
	if stats.Pools[0].Pity != 2 || !stats.Pools[0].PityIsFloor {
		t.Fatalf("无 5★ 时垫抽数应为窗口下界 2: %+v", stats.Pools[0])
	}
}

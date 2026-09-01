// Package gacha 实现抽卡域：凭证解析、官方接口拉取、增量合并、存储与统计。
// 数据结构与 Mac 版 aMC 完全互通（~/.amc/data/{UID}/gacha_data.json）。
package gacha

import (
	"net/url"
	"strings"
)

// 卡池类型与名称（与 models.py 的 POOL_TYPE_NAMES 一致）。
var PoolTypeNames = map[int]string{
	1: "角色活动唤取", 2: "武器活动唤取", 3: "角色常驻唤取", 4: "武器常驻唤取",
	5: "新手唤取", 6: "新手自选唤取", 7: "感恩定向唤取", 8: "角色新旅唤取",
	9: "武器新旅唤取", 10: "角色联动唤取", 11: "武器联动唤取", 12: "角色忆旅唤取",
	13: "武器忆旅唤取",
}

// AllPoolTypes 全部卡池类型，按 ID 升序。
var AllPoolTypes = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

const (
	APIBaseCN     = "https://gmserver-api.aki-game2.com/gacha"
	APIBaseGlobal = "https://gmserver-api.aki-game2.net/gacha"
)

// Credentials 从唤取记录 URL 解析出的请求凭证。
type Credentials struct {
	RecordID     string `json:"record_id"`
	PlayerID     string `json:"player_id"`
	ServerID     string `json:"server_id"`
	CardPoolID   string `json:"card_pool_id"`
	LanguageCode string `json:"language_code"`
	SvrArea      string `json:"svr_area"`
	IsGlobal     bool   `json:"is_global"`
}

// ParseCredentials 从唤取记录页 URL 提取凭证（移植自 GachaCredentials.from_url）。
func ParseCredentials(rawURL string) Credentials {
	normalized := strings.ReplaceAll(rawURL, "#", "")
	u, err := url.Parse(normalized)
	if err != nil {
		return Credentials{LanguageCode: "zh-Hans", SvrArea: "cn"}
	}
	q := u.Query()
	get := func(key string) string { return q.Get(key) }

	host := u.Hostname()
	isGlobal := strings.HasSuffix(host, ".net") || strings.Contains(host, "oversea")

	svrArea := get("svr_area")
	if svrArea == "" {
		if isGlobal {
			svrArea = "global"
		} else {
			svrArea = "cn"
		}
	}
	lang := get("lang")
	if lang == "" {
		lang = "zh-Hans"
	}
	return Credentials{
		RecordID:     get("record_id"),
		PlayerID:     get("player_id"),
		ServerID:     get("svr_id"),
		CardPoolID:   orDefault(get("resources_id"), get("gacha_id")),
		LanguageCode: lang,
		SvrArea:      svrArea,
		IsGlobal:     isGlobal,
	}
}

// APIBase 国服/国际服接口基地址。
func (c Credentials) APIBase() string {
	if c.IsGlobal {
		return APIBaseGlobal
	}
	return APIBaseCN
}

func orDefault(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

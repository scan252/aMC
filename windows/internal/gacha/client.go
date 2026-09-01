package gacha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Record 单条唤取记录（JSON 字段名与 Mac 版存储格式一致）。
type Record struct {
	CardPoolType string `json:"cardPoolType"`
	ResourceID   int    `json:"resourceId"`
	QualityLevel int    `json:"qualityLevel"`
	ResourceType string `json:"resourceType"`
	Name         string `json:"name"`
	Count        int    `json:"count"`
	Time         string `json:"time"`
}

// UniqueKey 记录去重键：time|name|qualityLevel（与 Mac 版 unique_key 一致）。
func (r Record) UniqueKey() string {
	return fmt.Sprintf("%s|%s|%d", r.Time, r.Name, r.QualityLevel)
}

const (
	requestInterval = 500 * time.Millisecond
	requestTimeout  = 15 * time.Second
	userAgent       = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

// apiResponse 官方接口响应结构。
type apiResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    []Record `json:"data"`
}

// ProgressFn 每完成一个卡池回调一次（index 从 1 开始）；poolName 为空时表示收尾。
type ProgressFn func(index, total int, poolName string, err error)

// FetchPoolRecords 拉取单个卡池的唤取记录。
func FetchPoolRecords(client *http.Client, creds Credentials, poolType int) ([]Record, error) {
	payload := map[string]any{
		"recordId":     creds.RecordID,
		"playerId":     creds.PlayerID,
		"serverId":     creds.ServerID,
		"cardPoolId":   creds.CardPoolID,
		"cardPoolType": poolType,
		"languageCode": creds.LanguageCode,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, creds.APIBase()+"/record/query", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	if client == nil {
		client = &http.Client{Timeout: requestTimeout}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("卡池 %d 请求失败: %w", poolType, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("卡池 %d HTTP %d", poolType, resp.StatusCode)
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result apiResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("卡池 %d 响应解析失败: %w", poolType, err)
	}
	switch {
	case result.Code == 0 || result.Code == 200:
	case result.Message == "成功" || result.Message == "success":
	default:
		return nil, fmt.Errorf("卡池 %d API 错误: code=%d, message=%s", poolType, result.Code, result.Message)
	}
	if result.Data == nil {
		return []Record{}, nil
	}
	return result.Data, nil
}

// FetchAllPools 顺序拉取全部 13 个卡池并合并为 pools 映射（键为卡池 ID 字符串）。
// 单个卡池失败不中断整体流程（与 Mac 版一致），错误通过 onProgress 上报。
func FetchAllPools(creds Credentials, onProgress ProgressFn) map[string][]Record {
	pools := make(map[string][]Record, len(AllPoolTypes))
	client := &http.Client{Timeout: requestTimeout}
	total := len(AllPoolTypes)

	for i, poolType := range AllPoolTypes {
		poolName := PoolTypeNames[poolType]
		if onProgress != nil {
			onProgress(i+1, total, poolName, nil)
		}
		records, err := FetchPoolRecords(client, creds, poolType)
		if err != nil {
			if onProgress != nil {
				onProgress(i+1, total, poolName, err)
			}
			records = []Record{}
		}
		pools[fmt.Sprint(poolType)] = records

		if i < total-1 {
			time.Sleep(requestInterval)
		}
	}
	return pools
}

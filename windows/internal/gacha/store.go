package gacha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// GachaData 落盘数据结构（与 Mac 版 GachaData.to_dict 完全同构）。
type GachaData struct {
	PlayerID  string              `json:"player_id"`
	SvrArea   string              `json:"svr_area"`
	FetchedAt string              `json:"fetched_at"`
	Pools     map[string][]Record `json:"pools"`
}

// DefaultDataDir 默认数据目录：~/.amc/data。
func DefaultDataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".amc", "data"), nil
}

// PlayerDataPath 指定 UID 的数据文件路径：{dataDir}/{uid}/gacha_data.json。
func PlayerDataPath(dataDir, playerID string) string {
	return filepath.Join(dataDir, playerID, "gacha_data.json")
}

// LoadGachaData 读取已有数据；文件不存在返回 nil 而非错误。
func LoadGachaData(path string) (*GachaData, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var data GachaData
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("解析 %s 失败: %w", path, err)
	}
	return &data, nil
}

// BackupExisting 在写入前把现有数据备份到同级 backup/ 目录。
func BackupExisting(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	backupDir := filepath.Join(filepath.Dir(path), "backup")
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return err
	}
	stamp := time.Now().Format("2006-01-02_150405")
	return os.WriteFile(filepath.Join(backupDir, "gacha_data_"+stamp+".json"), raw, 0o644)
}

// SaveGachaData 备份并写入数据（UTF-8、缩进 2、中文不转义，与 Mac 版格式对齐）。
func SaveGachaData(path string, data *GachaData) error {
	if err := BackupExisting(path); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
}

// ListPlayerIDs 列出数据目录下所有已存在数据的 UID。
func ListPlayerIDs(dataDir string) []string {
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return nil
	}
	var ids []string
	for _, e := range entries {
		if e.IsDir() {
			if _, err := os.Stat(filepath.Join(dataDir, e.Name(), "gacha_data.json")); err == nil {
				ids = append(ids, e.Name())
			}
		}
	}
	return ids
}


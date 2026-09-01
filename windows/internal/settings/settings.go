// Package settings 应用设置服务：~/.amc/settings.json。
package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Settings 应用设置。
type Settings struct {
	Autostart  bool   `json:"autostart"`   // 开机自启
	SignInAuto bool   `json:"signInAuto"`  // 每日自动签到
	SignInHour int    `json:"signInHour"`  // 自动签到时刻（0-23）
	WaveNotify bool   `json:"waveNotify"`  // 波片回满提醒
	LogPath    string `json:"logPath"`     // 手动指定的客户端日志路径
	Language   string `json:"language"`    // 界面语言（预留）
}

// Defaults 默认值。
func Defaults() *Settings {
	return &Settings{
		Autostart:  false,
		SignInAuto: true,
		SignInHour: 8,
		WaveNotify: true,
		Language:   "zh-CN",
	}
}

// Path 设置文件路径。
func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".amc", "settings.json"), nil
}

// Load 读取设置；文件不存在时返回默认值。
func Load() (*Settings, error) {
	s := Defaults()
	path, err := Path()
	if err != nil {
		return s, err
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return s, err
	}
	_ = json.Unmarshal(raw, s)
	if s.SignInHour < 0 || s.SignInHour > 23 {
		s.SignInHour = 8
	}
	return s, nil
}

// Save 保存设置。
func Save(s *Settings) error {
	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0o600)
}

// DataDir 当前数据目录（供设置页展示）。
func DataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".amc", "data"), nil
}

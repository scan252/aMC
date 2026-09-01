package kurobbs

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// NewAccount 组装账号并生成设备标识。
func NewAccount(token, userID, name, phone string) *Account {
	return &Account{
		Token:    token,
		UserID:   userID,
		Name:     name,
		Phone:    phone,
		DevCode:  randomHex(16),
		Did:      randomHex(16),
		LoggedIn: time.Now(),
	}
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Store 凭证本地存储（~/.amc/kurobbs.json）。
// 说明：底座阶段以用户目录 ACL 保护；后续版本计划接入 DPAPI 加密。
type Store struct {
	Path string
}

// DefaultStore 默认凭证存储位置。
func DefaultStore() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return &Store{Path: filepath.Join(home, ".amc", "kurobbs.json")}, nil
}

// Load 读取已保存账号；未登录返回 nil。
func (s *Store) Load() (*Account, error) {
	raw, err := os.ReadFile(s.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var acc Account
	if err := json.Unmarshal(raw, &acc); err != nil {
		return nil, fmt.Errorf("凭证文件损坏: %w", err)
	}
	return &acc, nil
}

// Save 保存账号凭证（目录权限 0700，文件 0600）。
func (s *Store) Save(acc *Account) error {
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o700); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(acc, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.Path, raw, 0o600)
}

// Clear 清除已保存凭证（退出登录）。
func (s *Store) Clear() error {
	err := os.Remove(s.Path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

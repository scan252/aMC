package main

import (
	"github.com/scan252/aMC/windows/internal/settings"
)

// SettingsService 设置域对前端的入口。
type SettingsService struct{}

// GetSettings 返回当前设置与数据目录。
type SettingsWithPaths struct {
	*settings.Settings
	DataDir   string `json:"dataDir"`
	AutostartOn bool `json:"autostartOn"`
}

func (s *SettingsService) Get() (*SettingsWithPaths, error) {
	st, err := settings.Load()
	if err != nil {
		return nil, err
	}
	dir, _ := settings.DataDir()
	return &SettingsWithPaths{
		Settings:    st,
		DataDir:     dir,
		AutostartOn: GetAutostart(),
	}, nil
}

func (s *SettingsService) Save(st *settings.Settings) error {
	if err := settings.Save(st); err != nil {
		return err
	}
	return SetAutostart(st.Autostart)
}

package main

import (
	"golang.org/x/sys/windows/registry"
)

// autostartPath HKCU 的开机自启 Run 键。
const autostartKey = `Software\Microsoft\Windows\CurrentVersion\Run`
const autostartValue = "aMC Suite"

// SetAutostart 写入/移除开机自启（仅 HKCU，无需管理员权限）。
func SetAutostart(enable bool) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, autostartKey, registry.SET_VALUE|registry.QUERY_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	if !enable {
		_, _, err := key.GetStringValue(autostartValue)
		if err == registry.ErrNotExist {
			return nil
		}
		return key.DeleteValue(autostartValue)
	}
	exe, err := osExecutable()
	if err != nil {
		return err
	}
	return key.SetStringValue(autostartValue, `"`+exe+`"`)
}

// GetAutostart 查询当前自启状态。
func GetAutostart() bool {
	key, err := registry.OpenKey(registry.CURRENT_USER, autostartKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()
	_, _, err = key.GetStringValue(autostartValue)
	return err == nil
}

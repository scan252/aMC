package wulog

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// pathRe 匹配配置文本中的 Windows 盘符路径（JSON 转义或原生反斜杠）。
var pathRe = regexp.MustCompile(`[A-Za-z]:[\\/][^"\r\n]*`)

// 常见安装根候选（相对各盘符根目录）。游戏日志固定位于：
//   <安装根>/Wuthering Waves Game/Client/Saved/Logs/Client.log
// 或个别精简安装为 <安装根>/Client/Saved/Logs/Client.log
var (
	driveRootCandidates = []string{
		"", "Games", "Game", "Games\\Steamlibrary", "Program Files", "Program Files (x86)",
		"WeGameGames", "KuroGames",
	}
	logRelPath = "Client\\Saved\\Logs\\Client.log"
)

// Discover 自动发现 Windows 版鸣潮客户端日志。
// 顺序：注册表卸载项 → KR 启动器配置目录 → 常见盘符路径扫描。
// 返回按 mtime 新到旧排序的结果；找不到返回空列表。
func Discover() ([]LogInfo, error) {
	seen := map[string]struct{}{}
	var found []LogInfo

	add := func(paths ...string) {
		for _, p := range paths {
			if p == "" {
				continue
			}
			if _, ok := seen[p]; ok {
				continue
			}
			if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
				seen[p] = struct{}{}
				found = append(found, LogInfo{Path: p, Label: "Windows 客户端"})
			}
		}
	}

	add(registryCandidates()...)
	add(launcherCandidates()...)
	add(driveScanCandidates()...)

	sort.SliceStable(found, func(i, j int) bool {
		mi, _ := os.Stat(found[i].Path)
		mj, _ := os.Stat(found[j].Path)
		return mi.ModTime().After(mj.ModTime())
	})
	return found, nil
}

// registryCandidates 从注册表卸载项中查找鸣潮安装目录（HKLM/HKCU × 64/32 位视图）。
func registryCandidates() []string {
	var out []string
	roots := []struct {
		key  registry.Key
		path string
	}{
		{registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`},
		{registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`},
		{registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`},
	}
	displayNames := []string{"wuthering waves", "鸣潮"}

	for _, root := range roots {
		uninst, err := registry.OpenKey(root.key, root.path, registry.ENUMERATE_SUB_KEYS)
		if err != nil {
			continue
		}
		names, _ := uninst.ReadSubKeyNames(-1)
		uninst.Close()

		for _, name := range names {
			k, err := registry.OpenKey(root.key, root.path+`\`+name, registry.QUERY_VALUE)
			if err != nil {
				continue
			}
			display, _, _ := k.GetStringValue("DisplayName")
			lower := strings.ToLower(display)
			matched := false
			for _, dn := range displayNames {
				if strings.Contains(lower, dn) {
					matched = true
					break
				}
			}
			if !matched {
				k.Close()
				continue
			}
			installLoc, _, _ := k.GetStringValue("InstallLocation")
			k.Close()
			if installLoc != "" {
				out = append(out, logCandidatesFor(strings.TrimSpace(installLoc))...)
			}
		}
	}
	return out
}

// launcherCandidates 探测 KR 启动器配置中记录的游戏路径。
func launcherCandidates() []string {
	var out []string
	appdata, err := os.UserConfigDir()
	if err != nil {
		return out
	}
	launcherDirs := []string{
		filepath.Join(appdata, "KRSDK", "KRSDKExe"),
		filepath.Join(appdata, "krStarter"),
		filepath.Join(appdata, "KRLauncher"),
	}
	for _, dir := range launcherDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || (!strings.HasSuffix(e.Name(), ".json") && !strings.HasSuffix(e.Name(), ".config")) {
				continue
			}
			raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				continue
			}
			out = append(out, extractPathsFromConfig(string(raw))...)
		}
	}
	return out
}

// extractPathsFromConfig 从启动器配置文本里提取疑似游戏安装路径并转换为日志路径。
func extractPathsFromConfig(content string) []string {
	var out []string
	matches := pathRe.FindAllString(content, -1)
	for _, m := range matches {
		m = strings.ReplaceAll(m, `\\`, `\`)
		m = strings.ReplaceAll(m, "/", `\`)
		if !strings.Contains(strings.ToLower(m), "wuthering") &&
			!strings.Contains(m, "鸣潮") && !strings.Contains(strings.ToLower(m), "client") {
			continue
		}
		out = append(out, logCandidatesFor(m)...)
	}
	return out
}

// driveScanCandidates 扫描各固定盘符的常见安装位置。
func driveScanCandidates() []string {
	var out []string
	for c := 'C'; c <= 'Z'; c++ {
		drive := string(c) + `:\`
		if _, err := os.Stat(drive); err != nil {
			continue
		}
		for _, sub := range driveRootCandidates {
			root := drive
			if sub != "" {
				root = filepath.Join(drive, sub)
			}
			if _, err := os.Stat(root); err != nil {
				continue
			}
			if candidates := logCandidatesFor(root); len(candidates) > 0 {
				out = append(out, candidates...)
			}
		}
	}
	return out
}

// logCandidatesFor 对给定目录返回其中可能存在的日志文件路径（只做 Stat 级检查）。
// 兼容三种指向：目录即 Client 目录 / 目录为游戏根 / 目录为游戏根上一层。
func logCandidatesFor(dir string) []string {
	dir = strings.TrimRight(strings.TrimSpace(dir), `\/`)
	if dir == "" {
		return nil
	}
	candidates := []string{
		filepath.Join(dir, "Saved", "Logs", "Client.log"),
		filepath.Join(dir, "Client", "Saved", "Logs", "Client.log"),
		filepath.Join(dir, "Wuthering Waves Game", "Client", "Saved", "Logs", "Client.log"),
	}
	var out []string
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			out = append(out, c)
		}
	}
	return out
}

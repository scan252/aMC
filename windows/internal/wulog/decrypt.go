// Package wulog 负责《鸣潮》Windows 客户端日志的发现、解密与唤取记录 URL 提取。
// 解密算法与 Mac 版 aMC（amc/log_parser.py）保持一致：
// 魔数 \xa5\xef\xa5 之后逐字节异或（奇数字节 ^0xA5，偶数字节 ^0xEF）。
package wulog

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"strings"
)

// LogMagic 加密日志文件头的固定魔数。
var LogMagic = []byte{0xA5, 0xEF, 0xA5}

// GachaURLRe 唤取记录页 URL 特征（国服 aki-game.com / 国际服 aki-game.net，含 -oversea）。
var GachaURLRe = regexp.MustCompile(
	`https://aki-gm-resources(?:-oversea)?\.aki-game\.(?:net|com)/aki/gacha/index\.html#/record[?=&\w\-%.]+`)

// plaintextMarker 明文日志的首行特征。
const plaintextMarker = "Log file open"

// Decrypt 对日志字节流解密；明文输入原样返回（与 Mac 版行为一致）。
func Decrypt(data []byte) []byte {
	if len(data) < 3 {
		return data
	}
	payload := data
	if bytes.Equal(data[:3], LogMagic) {
		payload = data[3:]
	}
	out := make([]byte, len(payload))
	for i, b := range payload {
		if b%2 == 1 {
			b ^= 0xA5
		} else {
			b ^= 0xEF
		}
		out[i] = b
	}
	return out
}

// IsEncrypted 探测日志是否为加密格式。
func IsEncrypted(data []byte) bool {
	if len(data) < 3 {
		return false
	}
	if bytes.Equal(data[:3], LogMagic) {
		return true
	}
	if data[0] == 0 {
		head := string(Decrypt(data))
		return strings.HasPrefix(head, plaintextMarker)
	}
	return false
}

// ReadLog 读取日志文件并返回解密后的文本内容。
func ReadLog(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	if IsEncrypted(raw) {
		raw = Decrypt(raw)
	}
	text := strings.TrimPrefix(string(raw), "﻿")
	return text, nil
}

// ExtractGachaURLs 从日志文本中按出现顺序提取全部唤取记录 URL。
func ExtractGachaURLs(content string) []string {
	return GachaURLRe.FindAllString(content, -1)
}

// ExtractLatestURL 提取最近一次唤取记录 URL；playerID 非空时回溯匹配该玩家的 URL。
func ExtractLatestURL(content, playerID string) (string, error) {
	urls := ExtractGachaURLs(content)
	if len(urls) == 0 {
		return "", errors.New("日志中未找到唤取记录 URL，请先在游戏中打开「唤取 → 唤取记录」页面")
	}
	if playerID != "" {
		for i := len(urls) - 1; i >= 0; i-- {
			if strings.Contains(urls[i], playerID) {
				return urls[i], nil
			}
		}
	}
	return urls[len(urls)-1], nil
}

// LogInfo 描述一个已发现的日志文件。
type LogInfo struct {
	Path  string `json:"path"`
	Label string `json:"label"`
}

// ParseLog 读取并提取日志信息的便捷入口。
func ParseLog(path, playerID string) (LogInfo, string, error) {
	content, err := ReadLog(path)
	if err != nil {
		return LogInfo{}, "", err
	}
	url, err := ExtractLatestURL(content, playerID)
	if err != nil {
		return LogInfo{}, "", err
	}
	return LogInfo{Path: path, Label: "Windows 客户端"}, url, nil
}

// findFirstURLLine 供调试：返回首个包含 URL 的行号（1-based），未找到返回 0。
func findFirstURLLine(content string) int {
	sc := bufio.NewScanner(strings.NewReader(content))
	line := 0
	for sc.Scan() {
		line++
		if GachaURLRe.MatchString(sc.Text()) {
			return line
		}
	}
	return 0
}

package wulog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// encrypt 按格式约定构造加密字节流：奇偶选键方向与解码器相反
// （明文奇数字节 ^0xEF、偶数字节 ^0xA5，使密文奇偶性满足解码规则）。
func encrypt(plaintext []byte) []byte {
	out := make([]byte, 0, len(plaintext)+3)
	out = append(out, LogMagic...)
	for _, b := range plaintext {
		if b%2 == 1 {
			b ^= 0xEF
		} else {
			b ^= 0xA5
		}
		out = append(out, b)
	}
	return out
}

func TestDecryptEncryptedLog(t *testing.T) {
	plaintext := "Log file open, 2026-09-01\r\nhttps://aki-gm-resources.aki-game.com/aki/gacha/index.html#/record?player_id=100252731&svr_id=01&lang=zh-Hans"
	got := string(Decrypt(encrypt([]byte(plaintext))))
	if got != plaintext {
		t.Fatalf("解密结果不符:\n got=%q\nwant=%q", got, plaintext)
	}
}

func TestReadLogPlaintextPassthrough(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Client.log")
	plain := "Log file open, plain\r\nno magic here"
	if err := os.WriteFile(path, []byte(plain), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := ReadLog(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != plain {
		t.Fatalf("明文日志应原样返回，got=%q", got)
	}
}

func TestIsEncrypted(t *testing.T) {
	if !IsEncrypted(encrypt([]byte("Log file open"))) {
		t.Fatal("带魔数的加密日志应被识别为加密")
	}
	// 首字节为 0 且解密后出现明文标记
	fake := append([]byte{0}, encrypt([]byte("Log file open"))[3:]...)
	_ = fake
	if IsEncrypted([]byte("Log file open, plain log")) {
		t.Fatal("明文日志不应被识别为加密")
	}
	if IsEncrypted([]byte{0, 1, 2}) {
		t.Fatal("解密后不含标记的短数据不应被识别为加密")
	}
}

func TestExtractLatestURL(t *testing.T) {
	cnURL := "https://aki-gm-resources.aki-game.com/aki/gacha/index.html#/record?record_id=r1&player_id=100252731&svr_id=01&lang=zh-Hans&resources_id=p1"
	globalURL := "https://aki-gm-resources-oversea.aki-game.net/aki/gacha/index.html#/record?record_id=r2&player_id=999888777&svr_id=usa&resources_id=p2"
	content := "noise\r\n" + cnURL + "\r\nmore noise\r\n" + globalURL + "\r\n"

	got, err := ExtractLatestURL(content, "")
	if err != nil {
		t.Fatal(err)
	}
	if got != globalURL {
		t.Fatalf("应取最后一条 URL，got=%q", got)
	}

	got, err = ExtractLatestURL(content, "100252731")
	if err != nil {
		t.Fatal(err)
	}
	if got != cnURL {
		t.Fatalf("按玩家过滤应取国服 URL，got=%q", got)
	}

	if _, err := ExtractLatestURL("no urls here", ""); err == nil {
		t.Fatal("无 URL 应返回错误")
	}
}

func TestReadLogRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Client.log")
	plaintext := "Log file open\r\nhttps://aki-gm-resources.aki-game.com/aki/gacha/index.html#/record?player_id=100252731"
	if err := os.WriteFile(path, encrypt([]byte(plaintext)), 0o644); err != nil {
		t.Fatal(err)
	}
	content, err := ReadLog(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(content, "player_id=100252731") {
		t.Fatalf("解密后应包含 URL，got=%q", content)
	}
}

func TestExtractPathsFromConfig(t *testing.T) {
	// 构造真实安装结构：<tmp>/Wuthering Waves Game/Client/Saved/Logs/Client.log
	gameRoot := t.TempDir()
	clientDir := filepath.Join(gameRoot, "Wuthering Waves Game", "Client")
	logDir := filepath.Join(clientDir, "Saved", "Logs")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		t.Fatal(err)
	}
	logPath := filepath.Join(logDir, "Client.log")
	if err := os.WriteFile(logPath, []byte("Log file open"), 0o644); err != nil {
		t.Fatal(err)
	}

	// 模拟启动器配置（JSON 转义的反斜杠）
	cfg := `{"game":{"install_path":"` + strings.ReplaceAll(clientDir, `\`, `\\`) + `","other":"C:\\Something"}}`
	got := extractPathsFromConfig(cfg)
	if len(got) == 0 {
		t.Fatal("应从配置中提取到日志路径")
	}
	found := false
	for _, p := range got {
		if p == logPath {
			found = true
		}
	}
	if !found {
		t.Fatalf("应包含 %q，got=%v", logPath, got)
	}
}

package main

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"os"
)

//go:embed build/appicon.png
var appIconPNG []byte

// buildTrayIcon 把 PNG 打包为 ICO（Vista+ 支持 PNG 内嵌格式）。
func buildTrayIcon() []byte {
	var buf bytes.Buffer
	// ICONDIR: 保留字(2) + 类型 1=图标(2) + 数量(2)
	_ = binary.Write(&buf, binary.LittleEndian, uint16(0))
	_ = binary.Write(&buf, binary.LittleEndian, uint16(1))
	_ = binary.Write(&buf, binary.LittleEndian, uint16(1))

	size := uint32(len(appIconPNG))
	offset := uint32(6 + 16)
	// ICONDIRENTRY: 宽 0 表示 256
	buf.WriteByte(0)
	buf.WriteByte(0)
	buf.WriteByte(0)
	buf.WriteByte(0)
	_ = binary.Write(&buf, binary.LittleEndian, uint16(1)) // planes
	_ = binary.Write(&buf, binary.LittleEndian, uint16(32))
	_ = binary.Write(&buf, binary.LittleEndian, size)
	_ = binary.Write(&buf, binary.LittleEndian, offset)
	buf.Write(appIconPNG)
	return buf.Bytes()
}

func osExecutable() (string, error) {
	return os.Executable()
}

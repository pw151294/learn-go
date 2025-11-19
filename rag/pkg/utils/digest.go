package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// Sha256Hex 计算输入文本的 SHA-256 哈希，并返回十六进制字符串
func Sha256Hex(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
)

// 生成短链接代码
func GenerateShortCode(str string) string {
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:])[:8] // 8字符短代码
}

func PathExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

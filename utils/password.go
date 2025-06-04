package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// GenerateRandomString 生成安全随机字符串
func GenerateRandomString(length int) (string, error) {
	if length < 1 {
		return "", errors.New("长度必须大于0")
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	encoded := base64.URLEncoding.EncodeToString(bytes)
	return strings.TrimRight(encoded, "=")[:length], nil
}

// HashPassword 安全密码哈希（带盐值）
func HashPassword(password string, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(password + salt))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

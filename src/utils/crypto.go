package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateUUID 生成UUID v4
func GenerateUUID() string {
	// 生成16字节的随机数�?	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		// 如果随机数生成失败，使用时间�?随机数作为备选方�?		return generateFallbackUUID()
	}

	// 设置UUID版本 (4)
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	// 设置UUID变体 (RFC4122)
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	// 格式化为标准UUID格式
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])
}

// generateFallbackUUID 备选UUID生成方案
func generateFallbackUUID() string {
	// 使用时间戳和随机数生成UUID
	timestamp := time.Now().UnixNano()
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)

	// 组合时间戳和随机�?	combined := fmt.Sprintf("%d-%s", timestamp, hex.EncodeToString(randomBytes))
	
	// 生成哈希
	hash := sha256.Sum256([]byte(combined))
	
	// 格式化为UUID
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		hash[0:4], hash[4:6], hash[6:8], hash[8:10], hash[10:16])
}

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		// 从字符集中随机选择字符
		randomByte := make([]byte, 1)
		rand.Read(randomByte)
		b[i] = charset[randomByte[0]%byte(len(charset))]
	}
	return string(b)
}

// GenerateSecureToken 生成安全的随机令牌（用于邮箱验证、密码重置等�?func GenerateSecureToken() string {
	// 生成32字节的随机数�?	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		// 如果失败，使用备选方�?		return GenerateRandomString(64)
	}
	
	// 编码为十六进制字符串
	return hex.EncodeToString(token)
}

// GenerateShortCode 生成短验证码�?位数字）
func GenerateShortCode() string {
	code := make([]byte, 6)
	rand.Read(code)
	
	// 转换为数字字符串
	result := ""
	for _, b := range code {
		result += fmt.Sprintf("%d", b%10)
	}
	return result
}

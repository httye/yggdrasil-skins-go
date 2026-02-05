package utils

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// GenerateUserUUID 根据邮箱生成用户UUID（使用UUID v5�?func GenerateUserUUID(email string) string {
	// 使用DNS命名空间和邮箱生成UUID v5
	userUUID := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(strings.ToLower(email)))
	return strings.ReplaceAll(userUUID.String(), "-", "")
}

// GenerateProfileUUID 生成角色UUID
// 为了兼容离线验证，使用角色名称生成UUID
func GenerateProfileUUID(playerName string) string {
	// 兼容离线验证的UUID生成方法
	// UUID.nameUUIDFromBytes(("OfflinePlayer:" + characterName).getBytes(StandardCharsets.UTF_8))
	data := fmt.Sprintf("OfflinePlayer:%s", playerName)
	hash := md5.Sum([]byte(data))
	
	// 将MD5哈希转换为UUID格式
	// 设置版本位（�?字节的高4位设�?�?	hash[6] = (hash[6] & 0x0f) | 0x30
	// 设置变体位（�?字节的高2位设�?0�?	hash[8] = (hash[8] & 0x3f) | 0x80
	
	// 格式化为UUID字符串并移除连字�?	uuidStr := fmt.Sprintf("%08x%04x%04x%04x%012x",
		hash[0:4], hash[4:6], hash[6:8], hash[8:10], hash[10:16])
	
	return uuidStr
}

// GenerateRandomUUID 生成随机UUID（UUID v4�?func GenerateRandomUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// FormatUUID 格式化UUID（添加连字符�?func FormatUUID(uuidStr string) string {
	if len(uuidStr) != 32 {
		return uuidStr
	}
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		uuidStr[0:8], uuidStr[8:12], uuidStr[12:16], uuidStr[16:20], uuidStr[20:32])
}

// RemoveUUIDHyphens 移除UUID中的连字�?func RemoveUUIDHyphens(uuidStr string) string {
	return strings.ReplaceAll(uuidStr, "-", "")
}

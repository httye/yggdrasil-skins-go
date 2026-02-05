package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// ParseInt 安全地将字符串解析为整数
func ParseInt(s string) int {
	if s == "" {
		return 0
	}
	
	i, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0
	}
	return i
}

// ParseInt64 安全地将字符串解析为int64
func ParseInt64(s string) int64 {
	if s == "" {
		return 0
	}
	
	i, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// ParseFloat 安全地将字符串解析为浮点�?func ParseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0
	}
	return f
}

// ParseBool 安全地将字符串解析为布尔�?func ParseBool(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "true" || s == "1" || s == "yes" || s == "on"
}

// TruncateString 截断字符串到指定长度
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

// RemoveWhitespace 移除字符串中的所有空白字�?func RemoveWhitespace(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

// RemoveExtraWhitespace 移除多余的空白字符（保留单个空格�?func RemoveExtraWhitespace(s string) string {
	// 使用正则表达式替换多个空格为一个空�?	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(s, " "))
}

// SanitizeFilename 清理文件名，移除非法字符
func SanitizeFilename(filename string) string {
	// 移除非法字符
	re := regexp.MustCompile(`[<>:"/\\|?*]`)
	filename = re.ReplaceAllString(filename, "_")
	
	// 移除控制字符
	re = regexp.MustCompile(`[\x00-\x1f\x7f]`)
	filename = re.ReplaceAllString(filename, "")
	
	// 限制长度
	if len(filename) > 255 {
		filename = filename[:255]
	}
	
	return strings.TrimSpace(filename)
}

// ContainsString 检查字符串数组是否包含指定字符�?func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ContainsStringIgnoreCase 检查字符串数组是否包含指定字符串（忽略大小写）
func ContainsStringIgnoreCase(slice []string, item string) bool {
	itemLower := strings.ToLower(item)
	for _, s := range slice {
		if strings.ToLower(s) == itemLower {
			return true
		}
	}
	return false
}

// RemoveDuplicates 移除字符串数组中的重复项
func RemoveDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// SplitAndTrim 分割字符串并修剪每个部分
func SplitAndTrim(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	
	return result
}

// JoinNonEmpty 连接非空字符�?func JoinNonEmpty(sep string, parts ...string) string {
	nonEmpty := make([]string, 0, len(parts))
	
	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			nonEmpty = append(nonEmpty, part)
		}
	}
	
	if len(nonEmpty) == 0 {
		return ""
	}
	
	return strings.Join(nonEmpty, sep)
}

// MaskEmail 遮罩邮箱地址（保护隐私）
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}
	
	local := parts[0]
	domain := parts[1]
	
	// 遮罩本地部分
	if len(local) <= 3 {
		return "***@" + domain
	}
	
	maskedLocal := local[:1] + "***" + local[len(local)-1:]
	return maskedLocal + "@" + domain
}

// MaskPhoneNumber 遮罩手机号（保护隐私�?func MaskPhoneNumber(phone string) string {
	if len(phone) < 7 {
		return "***"
	}
	
	return phone[:3] + "****" + phone[len(phone)-4:]
}

// MaskQQNumber 遮罩QQ号码（保护隐私）
func MaskQQNumber(qq string) string {
	if len(qq) < 4 {
		return "***"
	}
	
	return qq[:1] + "****" + qq[len(qq)-2:]
}

// FormatBytes 格式化字节数（从utils/formatters.go移动过来�?func FormatBytes(bytes int64, decimals int) string {
	if bytes == 0 {
		return "0 Bytes"
	}
	
	k := int64(1024)
	sizes := []string{"Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	
	i := 0
	for bytes >= k && i < len(sizes)-1 {
		bytes /= k
		i++
	}
	
	if decimals < 0 {
		decimals = 0
	}
	
	return strconv.FormatInt(bytes, 10) + " " + sizes[i]
}

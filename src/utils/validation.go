package utils

import (
	"regexp"
	"strings"
)

var (
	// 邮箱验证正则表达式
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	
	// 用户名验证正则表达式（3-16位，字母数字下划线）
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,16}$`)
	
	// 游戏名验证正则表达式（Minecraft玩家名规则）
	playerNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,16}$`)
	
	// QQ号验证正则表达式
	qqNumberRegex = regexp.MustCompile(`^[1-9][0-9]{4,10}$`)
	
	// 密码强度验证正则表达式
	// 至少8位，包含大小写字母和数字
	strongPasswordRegex = regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9]).{8,}$`)
)

// IsValidEmail 验证邮箱格式
func IsValidEmail(email string) bool {
	if len(email) > 255 {
		return false
	}
	return emailRegex.MatchString(email)
}

// IsValidUsername 验证用户名格式
func IsValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 16 {
		return false
	}
	return usernameRegex.MatchString(username)
}

// IsValidPlayerName 验证游戏名格式（Minecraft玩家名规则）
func IsValidPlayerName(playerName string) bool {
	if len(playerName) < 3 || len(playerName) > 16 {
		return false
	}
	
	// Minecraft玩家名规则：
	// 1. 长度3-16个字符
	// 2. 只能包含字母、数字、下划线
	// 3. 不能以数字开头（某些服务器限制）
	// 4. 不能包含连续的下划线
	
	if !playerNameRegex.MatchString(playerName) {
		return false
	}
	
	// 检查是否以数字开头（可选限制）
	if len(playerName) > 0 && playerName[0] >= '0' && playerName[0] <= '9' {
		return false
	}
	
	// 检查是否包含连续的下划线
	if strings.Contains(playerName, "__") {
		return false
	}
	
	// 检查是否以下划线开头或结尾
	if strings.HasPrefix(playerName, "_") || strings.HasSuffix(playerName, "_") {
		return false
	}
	
	return true
}

// IsValidQQNumber 验证QQ号码格式
func IsValidQQNumber(qqNumber string) bool {
	if qqNumber == "" {
		return true // QQ号码是可选的
	}
	return qqNumberRegex.MatchString(qqNumber)
}

// IsValidPassword 验证密码格式
func IsValidPassword(password string) bool {
	if len(password) < 6 {
		return false
	}
	return true
}

// IsStrongPassword 验证强密码格式
func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	return strongPasswordRegex.MatchString(password)
}

// SanitizePlayerName 清理游戏名
func SanitizePlayerName(playerName string) string {
	// 移除前后空格
	playerName = strings.TrimSpace(playerName)
	
	// 转换为小写（Minecraft玩家名不区分大小写）
	playerName = strings.ToLower(playerName)
	
	// 移除非法字符
	playerName = regexp.MustCompile(`[^a-zA-Z0-9_]`).ReplaceAllString(playerName, "")
	
	// 限制长度
	if len(playerName) > 16 {
		playerName = playerName[:16]
	}
	
	return playerName
}

// SanitizeUsername 清理用户名
func SanitizeUsername(username string) string {
	// 移除前后空格
	username = strings.TrimSpace(username)
	
	// 限制长度
	if len(username) > 16 {
		username = username[:16]
	}
	
	return username
}

// SanitizeEmail 清理邮箱地址
func SanitizeEmail(email string) string {
	// 移除前后空格和换行符
	email = strings.TrimSpace(email)
	
	// 转换为小写
	email = strings.ToLower(email)
	
	return email
}

// ValidateRegistrationData 验证注册数据
func ValidateRegistrationData(email, username, password, playerName, qqNumber string) []string {
	var errors []string
	
	if !IsValidEmail(email) {
		errors = append(errors, "邮箱格式不正确")
	}
	
	if !IsValidUsername(username) {
		errors = append(errors, "用户名格式不正确（3-16位，只能包含字母、数字、下划线）")
	}
	
	if !IsValidPassword(password) {
		errors = append(errors, "密码长度至少6位")
	}
	
	if !IsValidPlayerName(playerName) {
		errors = append(errors, "游戏名格式不正确（3-16位，只能包含字母、数字、下划线，不能以数字开头）")
	}
	
	if !IsValidQQNumber(qqNumber) {
		errors = append(errors, "QQ号码格式不正确")
	}
	
	return errors
}

// CheckPasswordStrength 检查密码强度
func CheckPasswordStrength(password string) (score int, feedback string) {
	score = 0
	
	// 长度检查
	if len(password) >= 8 {
		score++
	}
	if len(password) >= 12 {
		score++
	}
	
	// 复杂度检查
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)
	
	if hasLower {
		score++
	}
	if hasUpper {
		score++
	}
	if hasDigit {
		score++
	}
	if hasSpecial {
		score++
	}
	
	// 生成反馈
	switch score {
	case 0, 1:
		feedback = "密码强度很弱，建议增加长度和复杂度"
	case 2, 3:
		feedback = "密码强度较弱，建议增加大小写字母、数字或特殊字符"
	case 4, 5:
		feedback = "密码强度中等，可以继续增强"
	case 6, 7:
		feedback = "密码强度较强"
	default:
		feedback = "密码强度很强"
	}
	
	return score, feedback
}
package utils

import (
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 哈希密码
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}
	
	// 使用bcrypt生成密码哈希，cost为10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	
	return string(hashedPassword), nil
}

// VerifyPassword 验证密码
func VerifyPassword(hashedPassword, password string) error {
	if hashedPassword == "" || password == "" {
		return errors.New("password cannot be empty")
	}
	
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// IsPasswordStrong 检查密码强度
func IsPasswordStrong(password string) (bool, []string) {
	var issues []string
	
	// 长度检查
	if len(password) < 8 {
		issues = append(issues, "Password must be at least 8 characters long")
	}
	
	// 包含小写字母
	hasLower := false
	for _, char := range password {
		if char >= 'a' && char <= 'z' {
			hasLower = true
			break
		}
	}
	if !hasLower {
		issues = append(issues, "Password must contain at least one lowercase letter")
	}
	
	// 包含大写字母
	hasUpper := false
	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		issues = append(issues, "Password must contain at least one uppercase letter")
	}
	
	// 包含数字
	hasDigit := false
	for _, char := range password {
		if char >= '0' && char <= '9' {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		issues = append(issues, "Password must contain at least one digit")
	}
	
	// 包含特殊字符（可选，但推荐）
	hasSpecial := false
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	for _, char := range password {
		for _, special := range specialChars {
			if char == special {
				hasSpecial = true
				break
			}
		}
		if hasSpecial {
			break
		}
	}
	if !hasSpecial && len(issues) == 0 {
		issues = append(issues, "Consider adding special characters for better security")
	}
	
	return len(issues) == 0, issues
}

// ValidatePasswordStrength 验证密码强度（返回分数和建议）
func ValidatePasswordStrength(password string) (score int, feedback string) {
	score = 0
	
	// 长度加分
	if len(password) >= 8 {
		score += 2
	}
	if len(password) >= 12 {
		score += 1
	}
	
	// 复杂度加分
	if regexp.MustCompile(`[a-z]`).MatchString(password) {
		score += 1
	}
	if regexp.MustCompile(`[A-Z]`).MatchString(password) {
		score += 1
	}
	if regexp.MustCompile(`[0-9]`).MatchString(password) {
		score += 1
	}
	if regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password) {
		score += 2
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

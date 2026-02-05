package services

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/NewNanCity/YggdrasilGo/src/models"
)

var (
	// ErrUserNotFound 用户未找到错误
	ErrUserNotFound = errors.New("user not found")
	// ErrUserAlreadyBanned 用户已被封禁错误
	ErrUserAlreadyBanned = errors.New("user is already banned")
	// ErrUserNotBanned 用户未被封禁错误
	ErrUserNotBanned = errors.New("user is not banned")
	// ErrInsufficientPrivileges 权限不足错误
	ErrInsufficientPrivileges = errors.New("insufficient privileges")
)

// UserBanService 用户封禁服务
type UserBanService struct {
	db *gorm.DB
}

// NewUserBanService 创建用户封禁服务
func NewUserBanService(db *gorm.DB) *UserBanService {
	return &UserBanService{db: db}
}

// BanUser 封禁用户
func (s *UserBanService) BanUser(targetUserUUID, adminUUID, reason string) error {
	// 获取目标用户信息
	var targetUser models.EnhancedUser
	if err := s.db.Where("uuid = ?", targetUserUUID).First(&targetUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to find target user: %w", err)
	}

	// 检查是否已被封禁
	if targetUser.IsBanned {
		return ErrUserAlreadyBanned
	}

	// 获取管理员信息
	var admin models.EnhancedUser
	if err := s.db.Where("uuid = ?", adminUUID).First(&admin).Error; err != nil {
		return fmt.Errorf("failed to find admin user: %w", err)
	}

	// 检查管理员权限
	if !admin.IsAdmin && !admin.HasPermission("ban_user") {
		return ErrInsufficientPrivileges
	}

	// 执行封禁操作
	now := time.Now()
	targetUser.IsBanned = true
	targetUser.BannedReason = reason
	targetUser.BannedAt = &now
	targetUser.BannedBy = adminUUID

	if err := s.db.Save(&targetUser).Error; err != nil {
		return fmt.Errorf("failed to ban user: %w", err)
	}

	// 记录管理员操作日志
	logEntry := models.AdminLog{
		AdminUUID:      adminUUID,
		Action:         "ban_user",
		TargetUserUUID: &targetUserUUID,
		Details: models.JSONMap{
			"reason":         reason,
			"previous_status": "active",
			"ban_time":       now.Format(time.RFC3339),
		},
	}
	if err := s.db.Create(&logEntry).Error; err != nil {
		// 记录日志失败，但不影响封禁操作
		fmt.Printf("Failed to log ban action: %v\n", err)
	}

	return nil
}

// UnbanUser 解封用户
func (s *UserBanService) UnbanUser(targetUserUUID, adminUUID string) error {
	// 获取目标用户信息
	var targetUser models.EnhancedUser
	if err := s.db.Where("uuid = ?", targetUserUUID).First(&targetUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to find target user: %w", err)
	}

	// 检查是否被封禁
	if !targetUser.IsBanned {
		return ErrUserNotBanned
	}

	// 获取管理员信息
	var admin models.EnhancedUser
	if err := s.db.Where("uuid = ?", adminUUID).First(&admin).Error; err != nil {
		return fmt.Errorf("failed to find admin user: %w", err)
	}

	// 检查管理员权限
	if !admin.IsAdmin && !admin.HasPermission("unban_user") {
		return ErrInsufficientPrivileges
	}

	// 保存封禁前的信息用于日志
	previousReason := targetUser.BannedReason
	previousBannedAt := targetUser.BannedAt

	// 执行解封操作
	targetUser.IsBanned = false
	targetUser.BannedReason = ""
	targetUser.BannedAt = nil
	targetUser.BannedBy = ""

	if err := s.db.Save(&targetUser).Error; err != nil {
		return fmt.Errorf("failed to unban user: %w", err)
	}

	// 记录管理员操作日志
	logEntry := models.AdminLog{
		AdminUUID:      adminUUID,
		Action:         "unban_user",
		TargetUserUUID: &targetUserUUID,
		Details: models.JSONMap{
			"previous_reason": previousReason,
			"previous_banned_at": previousBannedAt,
			"unban_time": time.Now().Format(time.RFC3339),
		},
	}
	if err := s.db.Create(&logEntry).Error; err != nil {
		// 记录日志失败，但不影响解封操作
		fmt.Printf("Failed to log unban action: %v\n", err)
	}

	return nil
}

// IsUserBanned 检查用户是否被封禁
func (s *UserBanService) IsUserBanned(userUUID string) (bool, error) {
	var user models.EnhancedUser
	err := s.db.Where("uuid = ?", userUUID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, ErrUserNotFound
		}
		return false, fmt.Errorf("failed to find user: %w", err)
	}
	return user.IsBanned, nil
}

// GetBannedUserInfo 获取被封禁用户的详细信息
func (s *UserBanService) GetBannedUserInfo(userUUID string) (*models.EnhancedUser, error) {
	var user models.EnhancedUser
	err := s.db.Where("uuid = ? AND is_banned = ?", userUUID, true).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotBanned
		}
		return nil, fmt.Errorf("failed to find banned user: %w", err)
	}
	return &user, nil
}

// GetBannedUsers 获取所有被封禁的用户列表
func (s *UserBanService) GetBannedUsers(limit, offset int) ([]models.EnhancedUser, int64, error) {
	var users []models.EnhancedUser
	var total int64

	// 获取总数
	err := s.db.Model(&models.EnhancedUser{}).Where("is_banned = ?", true).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count banned users: %w", err)
	}

	// 获取分页数据
	err = s.db.Where("is_banned = ?", true).
		Limit(limit).
		Offset(offset).
		Order("banned_at DESC").
		Find(&users).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get banned users: %w", err)
	}

	return users, total, nil
}

// GetBanHistory 获取用户的封禁历史（从操作日志中获取）
func (s *UserBanService) GetBanHistory(userUUID string, limit int) ([]models.AdminLog, error) {
	var logs []models.AdminLog
	err := s.db.Where("target_user_uuid = ? AND action IN (?, ?)", 
		userUUID, "ban_user", "unban_user").
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get ban history: %w", err)
	}
	return logs, nil
}

// ResetUserPassword 管理员重置用户密码
func (s *UserBanService) ResetUserPassword(targetUserUUID, adminUUID, newPassword string) error {
	// 获取目标用户信息
	var targetUser models.EnhancedUser
	if err := s.db.Where("uuid = ?", targetUserUUID).First(&targetUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to find target user: %w", err)
	}

	// 获取管理员信息
	var admin models.EnhancedUser
	if err := s.db.Where("uuid = ?", adminUUID).First(&admin).Error; err != nil {
		return fmt.Errorf("failed to find admin user: %w", err)
	}

	// 检查管理员权限
	if !admin.IsAdmin && !admin.HasPermission("reset_user_password") {
		return ErrInsufficientPrivileges
	}

	// 这里应该使用密码加密服务，简化处理
	targetUser.Password = newPassword // 实际应该加密

	if err := s.db.Save(&targetUser).Error; err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}

	// 记录管理员操作日志
	logEntry := models.AdminLog{
		AdminUUID:      adminUUID,
		Action:         "reset_user_password",
		TargetUserUUID: &targetUserUUID,
		Details: models.JSONMap{
			"reset_time": time.Now().Format(time.RFC3339),
		},
	}
	if err := s.db.Create(&logEntry).Error; err != nil {
		fmt.Printf("Failed to log password reset action: %v\n", err)
	}

	return nil
}

// UpdateUserMaxProfiles 更新用户角色数量限制
func (s *UserBanService) UpdateUserMaxProfiles(targetUserUUID, adminUUID string, maxProfiles int) error {
	// 获取目标用户信息
	var targetUser models.EnhancedUser
	if err := s.db.Where("uuid = ?", targetUserUUID).First(&targetUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to find target user: %w", err)
	}

	// 获取管理员信息
	var admin models.EnhancedUser
	if err := s.db.Where("uuid = ?", adminUUID).First(&admin).Error; err != nil {
		return fmt.Errorf("failed to find admin user: %w", err)
	}

	// 检查管理员权限
	if !admin.IsAdmin && !admin.HasPermission("update_user_limits") {
		return ErrInsufficientPrivileges
	}

	// 保存之前的值用于日志
	previousMax := targetUser.MaxProfiles

	// 更新角色数量限制
	targetUser.MaxProfiles = maxProfiles
	if err := s.db.Save(&targetUser).Error; err != nil {
		return fmt.Errorf("failed to update max profiles: %w", err)
	}

	// 记录管理员操作日志
	logEntry := models.AdminLog{
		AdminUUID:      adminUUID,
		Action:         "update_user_max_profiles",
		TargetUserUUID: &targetUserUUID,
		Details: models.JSONMap{
			"previous_max": previousMax,
			"new_max":      maxProfiles,
			"update_time":  time.Now().Format(time.RFC3339),
		},
	}
	if err := s.db.Create(&logEntry).Error; err != nil {
		fmt.Printf("Failed to log max profiles update action: %v\n", err)
	}

	return nil
}

// GetUserManagementStats 获取用户管理统计
func (s *UserBanService) GetUserManagementStats() (map[string]interface{}, error) {
	var totalUsers, bannedUsers, adminUsers int64
	var newestUser models.EnhancedUser

	// 获取总用户数
	err := s.db.Model(&models.EnhancedUser{}).Count(&totalUsers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count total users: %w", err)
	}

	// 获取被封禁用户数
	err = s.db.Model(&models.EnhancedUser{}).Where("is_banned = ?", true).Count(&bannedUsers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count banned users: %w", err)
	}

	// 获取管理员数
	err = s.db.Model(&models.EnhancedUser{}).Where("is_admin = ?", true).Count(&adminUsers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count admin users: %w", err)
	}

	// 获取最新注册用户
	err = s.db.Order("created_at DESC").First(&newestUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get newest user: %w", err)
	}

	// 获取今日注册用户数
	var todayUsers int64
	today := time.Now().Truncate(24 * time.Hour)
	err = s.db.Model(&models.EnhancedUser{}).Where("created_at >= ?", today).Count(&todayUsers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count today users: %w", err)
	}

	// 获取本周管理员操作数
	var thisWeekAdminActions int64
	weekAgo := time.Now().AddDate(0, 0, -7)
	err = s.db.Model(&models.AdminLog{}).Where("created_at >= ?", weekAgo).Count(&thisWeekAdminActions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count this week admin actions: %w", err)
	}

	stats := map[string]interface{}{
		"total_users":              totalUsers,
		"banned_users":             bannedUsers,
		"admin_users":              adminUsers,
		"today_new_users":          todayUsers,
		"this_week_admin_actions":  thisWeekAdminActions,
		"ban_rate":                 float64(bannedUsers) / float64(totalUsers) * 100,
		"newest_user":              newestUser.Username,
		"newest_user_created":      newestUser.CreatedAt,
	}

	return stats, nil
}
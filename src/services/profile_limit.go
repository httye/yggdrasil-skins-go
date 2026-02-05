package services

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/NewNanCity/YggdrasilGo/src/models"
)

var (
	// ErrProfileLimitReached 角色数量限制错误
	ErrProfileLimitReached = errors.New("profile limit reached")
	// ErrUserBanned 用户被封禁错误
	ErrUserBanned = errors.New("user is banned")
	// ErrProfileNotFound 角色未找到错误
	ErrProfileNotFound = errors.New("profile not found")
	// ErrProfileNameExists 角色名已存在错误
	ErrProfileNameExists = errors.New("profile name already exists")
	// ErrProfileNotOwned 角色不属于用户错误
	ErrProfileNotOwned = errors.New("profile does not belong to user")
)

// ProfileLimitService 角色限制服务
type ProfileLimitService struct {
	db *gorm.DB
}

// NewProfileLimitService 创建角色限制服务
func NewProfileLimitService(db *gorm.DB) *ProfileLimitService {
	return &ProfileLimitService{db: db}
}

// CanCreateProfile 检查用户是否可以创建新角色
func (s *ProfileLimitService) CanCreateProfile(userUUID string) (bool, int, int, error) {
	var user models.EnhancedUser
	if err := s.db.Where("uuid = ?", userUUID).First(&user).Error; err != nil {
		return false, 0, 0, fmt.Errorf("failed to find user: %w", err)
	}

	// 检查用户是否被封禁
	if user.IsBanned {
		return false, 0, user.MaxProfiles, ErrUserBanned
	}

	// 获取当前活跃角色数量
	var currentCount int64
	err := s.db.Model(&models.Profile{}).
		Where("user_uuid = ? AND is_active = ?", userUUID, true).
		Count(&currentCount).Error
	if err != nil {
		return false, 0, user.MaxProfiles, fmt.Errorf("failed to count profiles: %w", err)
	}

	// 检查是否达到限制
	if user.MaxProfiles == -1 {
		return true, int(currentCount), user.MaxProfiles, nil // 无限制
	}

	canCreate := int(currentCount) < user.MaxProfiles
	return canCreate, int(currentCount), user.MaxProfiles, nil
}

// CreateProfile 创建角色（带限制检查）
func (s *ProfileLimitService) CreateProfile(userUUID, profileName string) (*models.Profile, error) {
	// 检查是否可以创建角色
	canCreate, currentCount, maxAllowed, err := s.CanCreateProfile(userUUID)
	if err != nil {
		return nil, err
	}

	if !canCreate {
		if maxAllowed == -1 {
			return nil, fmt.Errorf("%w: user is banned", ErrUserBanned)
		}
		return nil, fmt.Errorf("%w: current %d, maximum %d", ErrProfileLimitReached, currentCount, maxAllowed)
	}

	// 检查角色名是否已存在
	var existingProfile models.Profile
	err = s.db.Where("name = ?", profileName).First(&existingProfile).Error
	if err == nil {
		return nil, ErrProfileNameExists
	}
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check profile name: %w", err)
	}

	// 生成角色UUID（这里应该使用UUID生成库）
	profileUUID := generateUUID()

	// 创建角色
	profile := &models.Profile{
		UUID:     profileUUID,
		Name:     profileName,
		UserUUID: userUUID,
		IsActive: true,
	}

	if err := s.db.Create(profile).Error; err != nil {
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}

	return profile, nil
}

// GetUserProfiles 获取用户的所有角色
func (s *ProfileLimitService) GetUserProfiles(userUUID string) ([]models.Profile, error) {
	var profiles []models.Profile
	err := s.db.Where("user_uuid = ?", userUUID).Find(&profiles).Error
	return profiles, err
}

// GetUserActiveProfiles 获取用户的活跃角色
func (s *ProfileLimitService) GetUserActiveProfiles(userUUID string) ([]models.Profile, error) {
	var profiles []models.Profile
	err := s.db.Where("user_uuid = ? AND is_active = ?", userUUID, true).Find(&profiles).Error
	return profiles, err
}

// GetProfileByUUID 通过UUID获取角色
func (s *ProfileLimitService) GetProfileByUUID(profileUUID string) (*models.Profile, error) {
	var profile models.Profile
	err := s.db.Where("uuid = ?", profileUUID).First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}
	return &profile, nil
}

// GetProfileByName 通过名称获取角色
func (s *ProfileLimitService) GetProfileByName(profileName string) (*models.Profile, error) {
	var profile models.Profile
	err := s.db.Where("name = ?", profileName).First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}
	return &profile, nil
}

// UpdateProfile 更新角色信息
func (s *ProfileLimitService) UpdateProfile(profileUUID string, updates map[string]interface{}) error {
	return s.db.Model(&models.Profile{}).Where("uuid = ?", profileUUID).Updates(updates).Error
}

// DeleteProfile 删除角色
func (s *ProfileLimitService) DeleteProfile(profileUUID string) error {
	return s.db.Where("uuid = ?", profileUUID).Delete(&models.Profile{}).Error
}

// DeactivateProfile 停用角色
func (s *ProfileLimitService) DeactivateProfile(profileUUID string) error {
	return s.UpdateProfile(profileUUID, map[string]interface{}{
		"is_active": false,
	})
}

// ActivateProfile 激活角色
func (s *ProfileLimitService) ActivateProfile(profileUUID string) error {
	// 获取角色信息
	profile, err := s.GetProfileByUUID(profileUUID)
	if err != nil {
		return err
	}

	// 检查用户是否被封禁
	var user models.EnhancedUser
	if err := s.db.Where("uuid = ?", profile.UserUUID).First(&user).Error; err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	if user.IsBanned {
		return ErrUserBanned
	}

	return s.UpdateProfile(profileUUID, map[string]interface{}{
		"is_active": true,
	})
}

// ValidateProfileOwnership 验证角色所有权
func (s *ProfileLimitService) ValidateProfileOwnership(profileUUID, userUUID string) error {
	var count int64
	err := s.db.Model(&models.Profile{}).
		Where("uuid = ? AND user_uuid = ?", profileUUID, userUUID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrProfileNotOwned
	}
	return nil
}

// GetProfileUsageStats 获取角色使用统计
func (s *ProfileLimitService) GetProfileUsageStats(userUUID string) (map[string]interface{}, error) {
	var user models.EnhancedUser
	if err := s.db.Where("uuid = ?", userUUID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	var totalCount, activeCount int64
	err := s.db.Model(&models.Profile{}).Where("user_uuid = ?", userUUID).Count(&totalCount).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count total profiles: %w", err)
	}

	err = s.db.Model(&models.Profile{}).Where("user_uuid = ? AND is_active = ?", userUUID, true).Count(&activeCount).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count active profiles: %w", err)
	}

	maxProfiles := user.MaxProfiles
	usagePercentage := 0.0
	if maxProfiles > 0 {
		usagePercentage = float64(activeCount) / float64(maxProfiles) * 100
	}

	return map[string]interface{}{
		"total_profiles":     totalCount,
		"active_profiles":    activeCount,
		"max_profiles":       maxProfiles,
		"usage_percentage":   usagePercentage,
		"can_create_more":    maxProfiles == -1 || int(activeCount) < maxProfiles,
		"is_banned":          user.IsBanned,
	}, nil
}

// BatchDeactivateProfiles 批量停用用户角色（用于封禁用户）
func (s *ProfileLimitService) BatchDeactivateProfiles(userUUID string) error {
	return s.db.Model(&models.Profile{}).
		Where("user_uuid = ?", userUUID).
		Update("is_active", false).Error
}

// generateUUID 生成UUID（简化版本，实际应该使用UUID库）
func generateUUID() string {
	// 这里应该使用github.com/google/uuid或其他UUID库
	// 为了简化，这里返回一个模拟的UUID
	return fmt.Sprintf("%d-%d-%d-%d-%d", 
		time.Now().Unix(), 
		time.Now().UnixNano()%1000, 
		time.Now().UnixNano()%10000,
		time.Now().UnixNano()%100000,
		time.Now().UnixNano()%1000000)
}
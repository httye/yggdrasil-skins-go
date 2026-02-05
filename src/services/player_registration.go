package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/httye/yggdrasil-skins-go/src/models"
	"github.com/httye/yggdrasil-skins-go/src/utils"
)

var (
	// ErrPlayerNameExists 游戏名已存在错误
	ErrPlayerNameExists = errors.New("player name already exists")
	// ErrInvalidPlayerName 无效游戏名错�?	ErrInvalidPlayerName = errors.New("invalid player name")
	// ErrPlayerVerificationFailed 游戏名验证失败错�?	ErrPlayerVerificationFailed = errors.New("player name verification failed")
	// ErrEmailNotVerified 邮箱未验证错�?	ErrEmailNotVerified = errors.New("email not verified")
	// ErrTermsNotAccepted 用户协议未接受错�?	ErrTermsNotAccepted = errors.New("terms not accepted")
)

// PlayerRegistrationService 游戏名注册服�?type PlayerRegistrationService struct {
	db                *gorm.DB
	yggdrasilAPIURL   string
	httpClient        *http.Client
}

// NewPlayerRegistrationService 创建游戏名注册服�?func NewPlayerRegistrationService(db *gorm.DB, yggdrasilAPIURL string) *PlayerRegistrationService {
	return &PlayerRegistrationService{
		db:              db,
		yggdrasilAPIURL: yggdrasilAPIURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// RegisterWithPlayerName 使用游戏名注册新用户
func (s *PlayerRegistrationService) RegisterWithPlayerName(request PlayerRegistrationRequest) (*models.EnhancedUser, error) {
	// 验证输入参数
	if err := s.validateRegistrationRequest(request); err != nil {
		return nil, err
	}

	// 检查游戏名是否已存�?	if err := s.checkPlayerNameAvailability(request.PlayerName); err != nil {
		return nil, err
	}

	// 验证邮箱格式
	if !utils.IsValidEmail(request.Email) {
		return nil, errors.New("invalid email format")
	}

	// 验证用户名格�?	if !utils.IsValidUsername(request.Username) {
		return nil, errors.New("invalid username format")
	}

	// 验证游戏名格�?	if !utils.IsValidPlayerName(request.PlayerName) {
		return nil, ErrInvalidPlayerName
	}

	// 验证游戏名和密码（通过Yggdrasil API�?	playerInfo, err := s.verifyPlayerCredentials(request.PlayerName, request.PlayerPassword)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPlayerVerificationFailed, err)
	}

	// 生成用户UUID
	userUUID := utils.GenerateUUID()

	// 创建用户记录
	user := &models.EnhancedUser{
		UUID:              userUUID,
		Email:             request.Email,
		Username:          request.Username,
		Password:          request.Password, // 应该已经加密
		PrimaryPlayerName: request.PlayerName,
		PlayerUUID:        playerInfo.UUID,
		QQNumber:          request.QQNumber,
		EmailVerified:     false, // 需要后续验�?		AgreedToTerms:     request.AgreedToTerms,
		RegistrationIP:    request.RegistrationIP,
		MaxProfiles:       5, // 默认限制
		IsAdmin:           false,
		PermissionGroupID: 1, // 默认权限�?	}

	// 保存用户到数据库
	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 记录注册日志
	s.logUserAction(user.UUID, "user_registered", models.JSONMap{
		"player_name": request.PlayerName,
		"player_uuid": playerInfo.UUID,
		"email":       request.Email,
		"username":    request.Username,
		"qq_number":   request.QQNumber,
	}, request.RegistrationIP, request.UserAgent)

	return user, nil
}

// PlayerRegistrationRequest 游戏名注册请�?type PlayerRegistrationRequest struct {
	Email             string `json:"email" binding:"required,email"`
	Username          string `json:"username" binding:"required,min=3,max=16,alphanum"`
	Password          string `json:"password" binding:"required,min=6"`
	PlayerName        string `json:"player_name" binding:"required,min=3,max=16"`
	PlayerPassword    string `json:"player_password" binding:"required"`
	QQNumber          string `json:"qq_number,omitempty"`
	AgreedToTerms     bool   `json:"agreed_to_terms" binding:"required,eq=true"`
	RegistrationIP    string `json:"-"`
	UserAgent         string `json:"-"`
}

// PlayerInfo 游戏玩家信息
type PlayerInfo struct {
	UUID      string `json:"id"`
	Name      string `json:"name"`
	AccessToken string `json:"accessToken"`
	ClientToken string `json:"clientToken"`
}

// verifyPlayerCredentials 通过Yggdrasil API验证游戏名凭�?func (s *PlayerRegistrationService) verifyPlayerCredentials(playerName, playerPassword string) (*PlayerInfo, error) {
	// 构建认证请求
	authRequest := map[string]interface{}{
		"username": playerName,
		"password": playerPassword,
		"agent": map[string]interface{}{
			"name":    "Minecraft",
			"version": 1,
		},
	}

	// 序列化请求数�?	requestBody, err := json.Marshal(authRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal auth request: %w", err)
	}

	// 发送认证请�?	url := s.yggdrasilAPIURL + "/authserver/authenticate"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 发送请�?	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send auth request: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状�?	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, fmt.Errorf("authentication failed with status %d", resp.StatusCode)
		}
		
		if errorMsg, ok := errorResp["errorMessage"].(string); ok {
			return nil, errors.New(errorMsg)
		}
		return nil, errors.New("authentication failed")
	}

	// 解析响应
	var authResponse struct {
		AccessToken   string `json:"accessToken"`
		ClientToken   string `json:"clientToken"`
		SelectedProfile struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"selectedProfile"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return nil, fmt.Errorf("failed to decode auth response: %w", err)
	}

	// 验证返回的游戏名是否匹配
	if authResponse.SelectedProfile.Name != playerName {
		return nil, errors.New("player name mismatch in authentication response")
	}

	return &PlayerInfo{
		UUID:        authResponse.SelectedProfile.ID,
		Name:        authResponse.SelectedProfile.Name,
		AccessToken: authResponse.AccessToken,
		ClientToken: authResponse.ClientToken,
	}, nil
}

// checkPlayerNameAvailability 检查游戏名是否可用
func (s *PlayerRegistrationService) checkPlayerNameAvailability(playerName string) error {
	var count int64
	err := s.db.Model(&models.EnhancedUser{}).Where("primary_player_name = ?", playerName).Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to check player name availability: %w", err)
	}
	if count > 0 {
		return ErrPlayerNameExists
	}
	return nil
}

// validateRegistrationRequest 验证注册请求
func (s *PlayerRegistrationService) validateRegistrationRequest(request PlayerRegistrationRequest) error {
	// 检查用户协议是否同�?	if !request.AgreedToTerms {
		return ErrTermsNotAccepted
	}

	// 检查邮箱是否已存在
	var emailCount int64
	err := s.db.Model(&models.EnhancedUser{}).Where("email = ?", request.Email).Count(&emailCount).Error
	if err != nil {
		return fmt.Errorf("failed to check email availability: %w", err)
	}
	if emailCount > 0 {
		return errors.New("email already exists")
	}

	// 检查用户名是否已存�?	var usernameCount int64
	err = s.db.Model(&models.EnhancedUser{}).Where("username = ?", request.Username).Count(&usernameCount).Error
	if err != nil {
		return fmt.Errorf("failed to check username availability: %w", err)
	}
	if usernameCount > 0 {
		return errors.New("username already exists")
	}

	return nil
}

// GetUserByPlayerName 通过游戏名获取用户信息（公开信息�?func (s *PlayerRegistrationService) GetUserByPlayerName(playerName string) (*PlayerUserInfo, error) {
	var user models.EnhancedUser
	err := s.db.Where("primary_player_name = ?", playerName).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("player not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// 返回公开信息（不包含敏感信息如密码、IP等）
	return &PlayerUserInfo{
		UUID:              user.UUID,
		Username:          user.Username,
		PrimaryPlayerName: user.PrimaryPlayerName,
		PlayerUUID:        user.PlayerUUID,
		QQNumber:          user.QQNumber,
		EmailVerified:     user.EmailVerified,
		CreatedAt:         user.CreatedAt,
	}, nil
}

// PlayerUserInfo 玩家用户信息（公开信息�?type PlayerUserInfo struct {
	UUID              string    `json:"uuid"`
	Username          string    `json:"username"`
	PrimaryPlayerName string    `json:"primary_player_name"`
	PlayerUUID        string    `json:"player_uuid"`
	QQNumber          string    `json:"qq_number,omitempty"`
	EmailVerified     bool      `json:"email_verified"`
	CreatedAt         time.Time `json:"created_at"`
}

// CheckPlayerNameExists 检查游戏名是否存在
func (s *PlayerRegistrationService) CheckPlayerNameExists(playerName string) (bool, error) {
	var count int64
	err := s.db.Model(&models.EnhancedUser{}).Where("primary_player_name = ?", playerName).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check player name existence: %w", err)
	}
	return count > 0, nil
}

// logUserAction 记录用户操作日志
func (s *PlayerRegistrationService) logUserAction(userUUID, action string, details models.JSONMap, ipAddress, userAgent string) {
	// 异步记录日志，不阻塞主流�?	go func() {
		if err := models.LogUserAction(s.db, userUUID, action, details, ipAddress, userAgent); err != nil {
			// 记录日志失败，可以在这里添加错误处理
			fmt.Printf("Failed to log user action: %v\n", err)
		}
	}()
}

// VerifyEmail 验证用户邮箱
func (s *PlayerRegistrationService) VerifyEmail(userUUID, token string) error {
	var user models.EnhancedUser
	if err := s.db.Where("uuid = ? AND email_verification_token = ?", userUUID, token).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("invalid verification token")
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	// 验证邮箱
	user.EmailVerified = true
	user.EmailVerificationToken = ""
	
	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	// 记录日志
	s.logUserAction(userUUID, "email_verified", models.JSONMap{
		"email": user.Email,
	}, "", "")

	return nil
}

// GenerateEmailVerificationToken 生成邮箱验证令牌
func (s *PlayerRegistrationService) GenerateEmailVerificationToken() string {
	return utils.GenerateRandomString(32)
}

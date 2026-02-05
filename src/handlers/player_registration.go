package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"yggdrasil-api-go/src/services"
	"yggdrasil-api-go/src/utils"
)

// PlayerRegistrationHandler 游戏名注册处理器
type PlayerRegistrationHandler struct {
	db                       *gorm.DB
	playerRegistrationService *services.PlayerRegistrationService
}

// NewPlayerRegistrationHandler 创建游戏名注册处理器
func NewPlayerRegistrationHandler(db *gorm.DB, yggdrasilAPIURL string) *PlayerRegistrationHandler {
	return &PlayerRegistrationHandler{
		db:                       db,
		playerRegistrationService: services.NewPlayerRegistrationService(db, yggdrasilAPIURL),
	}
}

// RegisterWithPlayerName 使用游戏名注册新用户
func (h *PlayerRegistrationHandler) RegisterWithPlayerName(c *gin.Context) {
	var request services.PlayerRegistrationRequest

	// 绑定请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	// 获取客户端信息
	request.RegistrationIP = c.ClientIP()
	request.UserAgent = c.Request.UserAgent()

	// 验证注册数据
	validationErrors := utils.ValidateRegistrationData(
		request.Email,
		request.Username,
		request.Password,
		request.PlayerName,
		request.QQNumber,
	)
	if len(validationErrors) > 0 {
		utils.RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", validationErrors[0])
		return
	}

	// 检查密码强度
	score, feedback := utils.CheckPasswordStrength(request.Password)
	if score < 3 {
		utils.RespondError(c, http.StatusBadRequest, "WEAK_PASSWORD", feedback)
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "PASSWORD_HASH_ERROR", "Failed to hash password")
		return
	}
	request.Password = hashedPassword

	// 执行注册
	user, err := h.playerRegistrationService.RegisterWithPlayerName(request)
	if err != nil {
		switch err {
		case services.ErrPlayerNameExists:
			utils.RespondError(c, http.StatusConflict, "PLAYER_NAME_EXISTS", "Player name already exists")
		case services.ErrInvalidPlayerName:
			utils.RespondError(c, http.StatusBadRequest, "INVALID_PLAYER_NAME", "Invalid player name format")
		case services.ErrPlayerVerificationFailed:
			utils.RespondError(c, http.StatusUnauthorized, "PLAYER_VERIFICATION_FAILED", "Player name or password is incorrect")
		case services.ErrTermsNotAccepted:
			utils.RespondError(c, http.StatusBadRequest, "TERMS_NOT_ACCEPTED", "You must agree to the terms of service")
		default:
			utils.RespondError(c, http.StatusInternalServerError, "REGISTRATION_FAILED", err.Error())
		}
		return
	}

	// 返回成功响应（不包含敏感信息）
	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
		"user": gin.H{
			"uuid":                user.UUID,
			"email":               user.Email,
			"username":            user.Username,
			"primary_player_name": user.PrimaryPlayerName,
			"player_uuid":         user.PlayerUUID,
			"qq_number":           user.QQNumber,
			"email_verified":      user.EmailVerified,
			"created_at":          user.CreatedAt,
		},
	})
}

// CheckPlayerNameAvailability 检查游戏名是否可用
func (h *PlayerRegistrationHandler) CheckPlayerNameAvailability(c *gin.Context) {
	playerName := c.Param("playerName")

	// 验证游戏名格式
	if !utils.IsValidPlayerName(playerName) {
		utils.RespondError(c, http.StatusBadRequest, "INVALID_PLAYER_NAME", "Invalid player name format")
		return
	}

	exists, err := h.playerRegistrationService.CheckPlayerNameExists(playerName)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to check player name availability")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"available": !exists,
		"exists":    exists,
	})
}

// GetUserByPlayerName 通过游戏名获取用户信息
func (h *PlayerRegistrationHandler) GetUserByPlayerName(c *gin.Context) {
	playerName := c.Param("playerName")

	// 验证游戏名格式
	if !utils.IsValidPlayerName(playerName) {
		utils.RespondError(c, http.StatusBadRequest, "INVALID_PLAYER_NAME", "Invalid player name format")
		return
	}

	userInfo, err := h.playerRegistrationService.GetUserByPlayerName(playerName)
	if err != nil {
		if err.Error() == "player not found" {
			utils.RespondError(c, http.StatusNotFound, "PLAYER_NOT_FOUND", "Player not found")
		} else {
			utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to find player")
		}
		return
	}

	c.JSON(http.StatusOK, userInfo)
}

// GetUserLogs 获取用户操作日志
func (h *PlayerRegistrationHandler) GetUserLogs(c *gin.Context) {
	userUUID := c.GetString("user_uuid")
	if userUUID == "" {
		utils.RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	limit := 50 // 默认限制50条
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit := utils.ParseInt(limitStr); parsedLimit > 0 && parsedLimit <= 200 {
			limit = parsedLimit
		}
	}

	logs, err := models.GetUserLogs(h.db, userUUID, limit)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to fetch user logs")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"count": len(logs),
	})
}

// VerifyEmail 验证用户邮箱
func (h *PlayerRegistrationHandler) VerifyEmail(c *gin.Context) {
	var request struct {
		UserUUID string `json:"user_uuid" binding:"required"`
		Token    string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := h.playerRegistrationService.VerifyEmail(request.UserUUID, request.Token); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "VERIFICATION_FAILED", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully",
	})
}

// SendEmailVerification 发送邮箱验证邮件
func (h *PlayerRegistrationHandler) SendEmailVerification(c *gin.Context) {
	userUUID := c.GetString("user_uuid")
	if userUUID == "" {
		utils.RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	// 获取用户信息
	var user models.EnhancedUser
	if err := h.db.Where("uuid = ?", userUUID).First(&user).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
		return
	}

	// 检查是否已验证
	if user.EmailVerified {
		utils.RespondError(c, http.StatusBadRequest, "EMAIL_ALREADY_VERIFIED", "Email already verified")
		return
	}

	// 生成验证令牌
	token := h.playerRegistrationService.GenerateEmailVerificationToken()
	user.EmailVerificationToken = token

	if err := h.db.Save(&user).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to generate verification token")
		return
	}

	// 这里应该发送邮件，简化处理
	// TODO: 集成邮件发送服务

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification email sent",
		"token":   token, // 开发环境下返回令牌，生产环境应该通过邮件发送
	})
}

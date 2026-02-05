package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/httye/yggdrasil-skins-go/src/models"
	"github.com/httye/yggdrasil-skins-go/src/services"
	"github.com/httye/yggdrasil-skins-go/src/utils"
)

// AdminHandler 后台管理处理�?type AdminHandler struct {
	db              *gorm.DB
	userBanService  *services.UserBanService
	profileService  *services.ProfileLimitService
}

// NewAdminHandler 创建后台管理处理�?func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{
		db:              db,
		userBanService:  services.NewUserBanService(db),
		profileService:  services.NewProfileLimitService(db),
	}
}

// GetUsers 获取用户列表
func (h *AdminHandler) GetUsers(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")
	isBanned := c.Query("is_banned")
	isAdmin := c.Query("is_admin")
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")

	// 构建查询
	query := h.db.Model(&models.EnhancedUser{})

	// 应用搜索条件
	if search != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 应用筛选条�?	if isBanned != "" {
		banned := isBanned == "true"
		query = query.Where("is_banned = ?", banned)
	}

	if isAdmin != "" {
		admin := isAdmin == "true"
		query = query.Where("is_admin = ?", admin)
	}

	// 应用排序
	if order == "desc" {
		query = query.Order(sort + " DESC")
	} else {
		query = query.Order(sort + " ASC")
	}

	// 执行分页查询
	var users []models.EnhancedUser
	var total int64

	if err := query.Count(&total).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to count users")
		return
	}

	offset := (page - 1) * pageSize
	if err := query.Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to fetch users")
		return
	}

	// 计算总页�?	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"data":        users,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

// GetUser 获取用户详细信息
func (h *AdminHandler) GetUser(c *gin.Context) {
	userUUID := c.Param("id")

	var user models.EnhancedUser
	if err := h.db.Where("uuid = ?", userUUID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to fetch user")
		return
	}

	// 获取用户完整信息
	userInfo, err := models.GetUserFullInfo(h.db, userUUID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to fetch user full info")
		return
	}

	c.JSON(http.StatusOK, userInfo)
}

// BanUser 封禁用户
func (h *AdminHandler) BanUser(c *gin.Context) {
	targetUserUUID := c.Param("id")
	adminUUID := c.GetString("user_uuid") // 从中间件获取当前管理员UUID

	var request struct {
		Reason string `json:"reason" binding:"required,min=1,max=500"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := h.userBanService.BanUser(targetUserUUID, adminUUID, request.Reason); err != nil {
		switch err {
		case services.ErrUserNotFound:
			utils.RespondError(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
		case services.ErrUserAlreadyBanned:
			utils.RespondError(c, http.StatusBadRequest, "USER_ALREADY_BANNED", "User is already banned")
		case services.ErrInsufficientPrivileges:
			utils.RespondError(c, http.StatusForbidden, "INSUFFICIENT_PRIVILEGES", "Insufficient privileges")
		default:
			utils.RespondError(c, http.StatusInternalServerError, "BAN_FAILED", "Failed to ban user")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User banned successfully",
		"reason":  request.Reason,
	})
}

// UnbanUser 解封用户
func (h *AdminHandler) UnbanUser(c *gin.Context) {
	targetUserUUID := c.Param("id")
	adminUUID := c.GetString("user_uuid") // 从中间件获取当前管理员UUID

	if err := h.userBanService.UnbanUser(targetUserUUID, adminUUID); err != nil {
		switch err {
		case services.ErrUserNotFound:
			utils.RespondError(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
		case services.ErrUserNotBanned:
			utils.RespondError(c, http.StatusBadRequest, "USER_NOT_BANNED", "User is not banned")
		case services.ErrInsufficientPrivileges:
			utils.RespondError(c, http.StatusForbidden, "INSUFFICIENT_PRIVILEGES", "Insufficient privileges")
		default:
			utils.RespondError(c, http.StatusInternalServerError, "UNBAN_FAILED", "Failed to unban user")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User unbanned successfully",
	})
}

// ResetUserPassword 重置用户密码
func (h *AdminHandler) ResetUserPassword(c *gin.Context) {
	targetUserUUID := c.Param("id")
	adminUUID := c.GetString("user_uuid")

	var request struct {
		NewPassword string `json:"new_password" binding:"required,min=6,max=100"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	// 这里应该使用密码加密服务
	// encryptedPassword := utils.HashPassword(request.NewPassword)

	if err := h.userBanService.ResetUserPassword(targetUserUUID, adminUUID, request.NewPassword); err != nil {
		switch err {
		case services.ErrUserNotFound:
			utils.RespondError(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
		case services.ErrInsufficientPrivileges:
			utils.RespondError(c, http.StatusForbidden, "INSUFFICIENT_PRIVILEGES", "Insufficient privileges")
		default:
			utils.RespondError(c, http.StatusInternalServerError, "PASSWORD_RESET_FAILED", "Failed to reset password")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
	})
}

// UpdateUserMaxProfiles 更新用户角色数量限制
func (h *AdminHandler) UpdateUserMaxProfiles(c *gin.Context) {
	targetUserUUID := c.Param("id")
	adminUUID := c.GetString("user_uuid")

	var request struct {
		MaxProfiles int `json:"max_profiles" binding:"required,min=-1"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := h.userBanService.UpdateUserMaxProfiles(targetUserUUID, adminUUID, request.MaxProfiles); err != nil {
		switch err {
		case services.ErrUserNotFound:
			utils.RespondError(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
		case services.ErrInsufficientPrivileges:
			utils.RespondError(c, http.StatusForbidden, "INSUFFICIENT_PRIVILEGES", "Insufficient privileges")
		default:
			utils.RespondError(c, http.StatusInternalServerError, "UPDATE_FAILED", "Failed to update max profiles")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Max profiles updated successfully",
		"max_profiles": request.MaxProfiles,
	})
}

// GetUserProfiles 获取用户的角色列�?func (h *AdminHandler) GetUserProfiles(c *gin.Context) {
	userUUID := c.Param("id")

	profiles, err := h.profileService.GetUserProfiles(userUUID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to fetch profiles")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profiles": profiles,
		"count":    len(profiles),
	})
}

// GetUserLogs 获取用户的操作日�?func (h *AdminHandler) GetUserLogs(c *gin.Context) {
	userUUID := c.Param("id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	logs, err := h.userBanService.GetBanHistory(userUUID, limit)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to fetch logs")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"count": len(logs),
	})
}

// GetStatistics 获取管理统计信息
func (h *AdminHandler) GetStatistics(c *gin.Context) {
	stats, err := h.userBanService.GetUserManagementStats()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to fetch statistics")
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetBannedUsers 获取被封禁用户列�?func (h *AdminHandler) GetBannedUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 计算偏移�?	offset := (page - 1) * pageSize

	users, total, err := h.userBanService.GetBannedUsers(pageSize, offset)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to fetch banned users")
		return
	}

	// 计算总页�?	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"data":        users,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

// AdminAuthMiddleware 管理员权限验证中间件
func (h *AdminHandler) AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID := c.GetString("user_uuid")
		if userUUID == "" {
			utils.RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			c.Abort()
			return
		}

		// 获取用户信息
		var user models.EnhancedUser
		if err := h.db.Where("uuid = ?", userUUID).First(&user).Error; err != nil {
			utils.RespondError(c, http.StatusUnauthorized, "USER_NOT_FOUND", "User not found")
			c.Abort()
			return
		}

		// 检查是否被封禁
		if user.IsBanned {
			utils.RespondError(c, http.StatusForbidden, "USER_BANNED", "User is banned")
			c.Abort()
			return
		}

		// 检查管理员权限
		if !user.IsAdmin {
			utils.RespondError(c, http.StatusForbidden, "INSUFFICIENT_PRIVILEGES", "Admin privileges required")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("admin_user", &user)
		c.Next()
	}
}

// GetCurrentAdmin 获取当前管理员信�?func (h *AdminHandler) GetCurrentAdmin(c *gin.Context) {
	admin := c.MustGet("admin_user").(*models.EnhancedUser)
	c.JSON(http.StatusOK, admin)
}

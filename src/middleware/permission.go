package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"yggdrasil-api-go/src/models"
	"yggdrasil-api-go/src/utils"
)

// PermissionMiddleware 权限验证中间件
type PermissionMiddleware struct {
	db *gorm.DB
}

// NewPermissionMiddleware 创建权限中间件
func NewPermissionMiddleware(db *gorm.DB) *PermissionMiddleware {
	return &PermissionMiddleware{db: db}
}

// RequirePermission 需要特定权限
func (pm *PermissionMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID := c.GetString("user_uuid")
		if userUUID == "" {
			utils.RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			c.Abort()
			return
		}

		// 获取用户信息
		var user models.EnhancedUser
		if err := pm.db.Where("uuid = ?", userUUID).First(&user).Error; err != nil {
			utils.RespondError(c, http.StatusUnauthorized, "USER_NOT_FOUND", "User not found")
			c.Abort()
			return
		}

		// 检查用户是否被封禁
		if user.IsBanned {
			utils.RespondError(c, http.StatusForbidden, "USER_BANNED", "User is banned")
			c.Abort()
			return
		}

		// 检查权限
		if !user.HasPermission(permission) {
			utils.RespondError(c, http.StatusForbidden, "INSUFFICIENT_PRIVILEGES", "Insufficient privileges")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin 需要管理员权限
func (pm *PermissionMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID := c.GetString("user_uuid")
		if userUUID == "" {
			utils.RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			c.Abort()
			return
		}

		// 获取用户信息
		var user models.EnhancedUser
		if err := pm.db.Where("uuid = ?", userUUID).First(&user).Error; err != nil {
			utils.RespondError(c, http.StatusUnauthorized, "USER_NOT_FOUND", "User not found")
			c.Abort()
			return
		}

		// 检查用户是否被封禁
		if user.IsBanned {
			utils.RespondError(c, http.StatusForbidden, "USER_BANNED", "User is banned")
			c.Abort()
			return
		}

		// 检查管理员权限
		if !user.IsAdmin {
			utils.RespondError(c, http.StatusForbidden, "ADMIN_REQUIRED", "Admin privileges required")
			c.Abort()
			return
		}

		// 将管理员信息存储到上下文中
		c.Set("admin_user", &user)
		c.Next()
	}
}

// RequirePermissionOrAdmin 需要特定权限或管理员权限
func (pm *PermissionMiddleware) RequirePermissionOrAdmin(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID := c.GetString("user_uuid")
		if userUUID == "" {
			utils.RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			c.Abort()
			return
		}

		// 获取用户信息
		var user models.EnhancedUser
		if err := pm.db.Where("uuid = ?", userUUID).First(&user).Error; err != nil {
			utils.RespondError(c, http.StatusUnauthorized, "USER_NOT_FOUND", "User not found")
			c.Abort()
			return
		}

		// 检查用户是否被封禁
		if user.IsBanned {
			utils.RespondError(c, http.StatusForbidden, "USER_BANNED", "User is banned")
			c.Abort()
			return
		}

		// 检查管理员权限或特定权限
		if !user.IsAdmin && !user.HasPermission(permission) {
			utils.RespondError(c, http.StatusForbidden, "INSUFFICIENT_PRIVILEGES", "Insufficient privileges")
			c.Abort()
			return
		}

		c.Next()
	}
}

// CheckUserBan 检查用户是否被封禁
func (pm *PermissionMiddleware) CheckUserBan() gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID := c.GetString("user_uuid")
		if userUUID == "" {
			c.Next() // 未登录用户继续执行
			return
		}

		// 获取用户信息
		var user models.EnhancedUser
		if err := pm.db.Where("uuid = ?", userUUID).First(&user).Error; err != nil {
			c.Next() // 用户不存在继续执行
			return
		}

		// 检查用户是否被封禁
		if user.IsBanned {
			utils.RespondError(c, http.StatusForbidden, "USER_BANNED", "User is banned")
			c.Abort()
			return
		}

		c.Next()
	}
}

// CheckProfileLimit 检查角色数量限制
func (pm *PermissionMiddleware) CheckProfileLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID := c.GetString("user_uuid")
		if userUUID == "" {
			utils.RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			c.Abort()
			return
		}

		// 获取用户信息
		var user models.EnhancedUser
		if err := pm.db.Where("uuid = ?", userUUID).First(&user).Error; err != nil {
			utils.RespondError(c, http.StatusUnauthorized, "USER_NOT_FOUND", "User not found")
			c.Abort()
			return
		}

		// 检查用户是否被封禁
		if user.IsBanned {
			utils.RespondError(c, http.StatusForbidden, "USER_BANNED", "User is banned")
			c.Abort()
			return
		}

		// 检查是否可以创建角色
		canCreate, _, _, err := user.CanCreateProfile(pm.db)
		if err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to check profile limit")
			c.Abort()
			return
		}

		if !canCreate {
			utils.RespondError(c, http.StatusBadRequest, "PROFILE_LIMIT_REACHED", "Profile limit reached")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByUser 基于用户的速率限制
func (pm *PermissionMiddleware) RateLimitByUser(maxRequests int, windowSeconds int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID := c.GetString("user_uuid")
		if userUUID == "" {
			c.Next() // 未登录用户不受限制
			return
		}

		// 这里应该实现基于Redis的速率限制
		// 简化实现，实际应该使用Redis或其他存储
		c.Next()
	}
}

// LogAdminAction 记录管理员操作
func (pm *PermissionMiddleware) LogAdminAction(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取管理员信息
		adminUser := c.MustGet("admin_user").(*models.EnhancedUser)
		
		// 获取目标用户UUID（如果存在）
		targetUserUUID := c.Param("id")
		
		// 获取请求详情
		details := map[string]interface{}{
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"user_agent":  c.Request.UserAgent(),
		}

		// 记录操作日志
		logEntry := models.AdminLog{
			AdminUUID:      adminUser.UUID,
			Action:         action,
			TargetUserUUID: &targetUserUUID,
			Details:        details,
			IPAddress:      c.ClientIP(),
			UserAgent:      c.Request.UserAgent(),
		}

		if err := pm.db.Create(&logEntry).Error; err != nil {
			// 记录日志失败，但不阻止操作
			// 可以在这里添加日志记录
		}

		c.Next()
	}
}

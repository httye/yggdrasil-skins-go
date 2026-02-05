package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/httye/yggdrasil-skins-go/src/handlers"
	"github.com/httye/yggdrasil-skins-go/src/middleware"
)

// SetupAdminRoutes 设置后台管理路由
func SetupAdminRoutes(router *gin.Engine, db *gorm.DB) {
	// 创建处理�?	adminHandler := handlers.NewAdminHandler(db)
	announcementHandler := handlers.NewAnnouncementHandler(db)

	// 后台管理API�?	admin := router.Group("/api/admin")
	{
		// 需要管理员权限验证
		admin.Use(adminHandler.AdminAuthMiddleware())

		// 用户管理
		admin.GET("/users", adminHandler.GetUsers)
		admin.GET("/users/:id", adminHandler.GetUser)
		admin.PUT("/users/:id/ban", adminHandler.BanUser)
		admin.PUT("/users/:id/unban", adminHandler.UnbanUser)
		admin.PUT("/users/:id/password", adminHandler.ResetUserPassword)
		admin.PUT("/users/:id/max-profiles", adminHandler.UpdateUserMaxProfiles)
		admin.GET("/users/:id/profiles", adminHandler.GetUserProfiles)
		admin.GET("/users/:id/logs", adminHandler.GetUserLogs)
		admin.GET("/banned-users", adminHandler.GetBannedUsers)

		// 公告管理
		admin.GET("/announcements", announcementHandler.GetAnnouncements)
		admin.POST("/announcements", announcementHandler.CreateAnnouncement)
		admin.GET("/announcements/:id", announcementHandler.GetAnnouncement)
		admin.PUT("/announcements/:id", announcementHandler.UpdateAnnouncement)
		admin.DELETE("/announcements/:id", announcementHandler.DeleteAnnouncement)
		admin.PUT("/announcements/:id/toggle", announcementHandler.ToggleAnnouncementStatus)

		// 系统统计
		admin.GET("/statistics", adminHandler.GetStatistics)
		admin.GET("/current-admin", adminHandler.GetCurrentAdmin)
	}

	// 公共公告API（不需要管理员权限�?	public := router.Group("/api")
	{
		// 获取有效公告
		public.GET("/announcements/active", announcementHandler.GetActiveAnnouncements)
		// 获取公告类型
		public.GET("/announcements/types", announcementHandler.GetAnnouncementTypes)
		// 获取目标用户�?		public.GET("/announcements/target-groups", announcementHandler.GetTargetGroups)
	}
}

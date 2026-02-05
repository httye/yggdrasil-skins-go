package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/httye/yggdrasil-skins-go/src/handlers"
	"github.com/httye/yggdrasil-skins-go/src/middleware"
)

// SetupPlayerRegistrationRoutes 设置游戏名注册路�?func SetupPlayerRegistrationRoutes(router *gin.Engine, db *gorm.DB, yggdrasilAPIURL string) {
	// 创建处理�?	playerHandler := handlers.NewPlayerRegistrationHandler(db, yggdrasilAPIURL)

	// 公开API组（不需要认证）
	public := router.Group("/api")
	{
		// 游戏名可用性检�?		public.GET("/players/check-name/:playerName", playerHandler.CheckPlayerNameAvailability)
		
		// 通过游戏名查询用户信息（公开信息�?		public.GET("/players/profile/:playerName", playerHandler.GetUserByPlayerName)
		
		// 增强注册（需要游戏名验证�?		public.POST("/auth/register-with-player", playerHandler.RegisterWithPlayerName)
	}

	// 需要认证的API�?	auth := router.Group("/api")
	auth.Use(middleware.JWTAuthMiddleware()) // JWT认证中间�?	{
		// 用户操作日志
		auth.GET("/users/logs", playerHandler.GetUserLogs)
		
		// 邮箱验证相关
		auth.POST("/auth/verify-email", playerHandler.VerifyEmail)
		auth.POST("/auth/send-email-verification", playerHandler.SendEmailVerification)
	}
}

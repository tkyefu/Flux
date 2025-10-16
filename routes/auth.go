package routes

import (
	"flux/database"
	"flux/handlers"
	"flux/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes 認証関連のルートを設定
func SetupAuthRoutes(router *gin.Engine) {
	handler := handlers.NewAuthHandler(database.DB)

	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)

		// 認証が必要なエンドポイント
		auth.GET("/me", middleware.AuthMiddleware(), handler.GetMe)
	}
}

package routes

import (
    "flux/handlers"
    "flux/mailer"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, mailer mailer.Mailer) {
    // ヘルスチェック
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // パスワードリセットハンドラー
    passwordResetHandler := handlers.NewPasswordResetHandler(db, mailer)

    // 認証関連のルート
    authGroup := r.Group("/api/auth")
    {
        // パスワードリセットリクエスト
        authGroup.POST("/forgot-password", passwordResetHandler.RequestReset)
        // パスワードリセット実行
        authGroup.POST("/reset-password", passwordResetHandler.ResetPassword)
    }
}
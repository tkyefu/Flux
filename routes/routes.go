package routes

import (
    "flux/handlers"
    "flux/mailer"
    "flux/middleware"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, mailer mailer.Mailer) {
    // ヘルスチェック
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    v1 := r.Group("/api/v1")
    {
        // 認証関連のルート
        authHandler := handlers.NewAuthHandler(db)
        auth := v1.Group("/auth")
        {
            auth.POST("/register", authHandler.Register)
            auth.POST("/login", authHandler.Login)
            auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetMe)

            // パスワードリセットハンドラー
            passwordResetHandler := handlers.NewPasswordResetHandler(db, mailer)
            auth.POST("/forgot-password", passwordResetHandler.RequestReset)
            auth.POST("/reset-password", passwordResetHandler.ResetPassword)
        }

        // tasks
        v1.GET("/tasks", handlers.GetTasks)
        v1.GET("/tasks/:id", handlers.GetTask)
        v1.POST("/tasks", middleware.AuthMiddleware(), handlers.CreateTask)
        v1.PUT("/tasks/:id", middleware.AuthMiddleware(), handlers.UpdateTask)
        v1.DELETE("/tasks/:id", middleware.AuthMiddleware(), handlers.DeleteTask)

        // users
        v1.GET("/users", handlers.GetUsers)
        v1.GET("/users/:id", handlers.GetUser)
        v1.POST("/users", handlers.CreateUser)
        v1.PUT("/users/:id", handlers.UpdateUser)
        v1.DELETE("/users/:id", handlers.DeleteUser)
        v1.GET("/users/:id/tasks", handlers.GetTasksByUser)
    }
}
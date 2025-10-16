package routes

import (
	"flux/database"
	"flux/handlers"
	"flux/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes (public)
		authHandler := handlers.NewAuthHandler(database.DB)
		v1.POST("/auth/register", authHandler.Register)
		v1.POST("/auth/login", authHandler.Login)
		v1.GET("/auth/me", middleware.AuthMiddleware(), authHandler.GetMe)

		// User-specific tasks (must be before generic user routes to avoid conflicts)
		v1.GET("/users/:id/tasks", middleware.AuthMiddleware(), handlers.GetTasksByUser)

		// User routes
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("", handlers.GetUsers)
			users.GET("/:id", handlers.GetUser)
			users.POST("", handlers.CreateUser) // 管理者のみ許可する場合は追加のミドルウェアが必要
			users.PUT("/:id", handlers.UpdateUser)
			users.DELETE("/:id", handlers.DeleteUser)
		}

		// Task routes
		tasks := v1.Group("/tasks")
		tasks.Use(middleware.AuthMiddleware())
		{
			tasks.GET("", handlers.GetTasks)
			tasks.GET("/:id", handlers.GetTask)
			tasks.POST("", handlers.CreateTask)
			tasks.PUT("/:id", handlers.UpdateTask)
			tasks.DELETE("/:id", handlers.DeleteTask)
		}
	}
}
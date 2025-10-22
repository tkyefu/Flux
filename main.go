package main

import (
    "log"
    "os"

    "flux/config"
    "flux/database"
    "flux/models"
    "flux/routes"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    // 環境変数の読み込み
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    // 設定の読み込み
    cfg := config.Load()

    // データベース接続
    db, err := database.Connect()
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // マイグレーション
    if len(os.Args) > 1 && os.Args[1] == "migrate" {
        if err := db.AutoMigrate(&models.User{}, &models.PasswordReset{}); err != nil {
            log.Fatalf("Failed to migrate database: %v", err)
        }
        log.Println("Migration completed successfully")
        return
    }

    // ルーターの設定
    r := gin.Default()
    routes.SetupRoutes(r, db, cfg.Mailer)

    // サーバー起動
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Server starting on port %s in %s mode\n", port, cfg.Env)
    log.Fatal(r.Run(":" + port))
}
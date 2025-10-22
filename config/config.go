package config

import (
    "flux/mailer"
    "os"
)

type Config struct {
    Env        string
    FrontendURL string
    Mailer     mailer.Mailer
}

func Load() *Config {
    env := getEnv("APP_ENV", "development")
    cfg := &Config{
        Env:        env,
        FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
    }

    // 環境に応じたメーラーを設定
    if env == "production" {
        cfg.Mailer = mailer.NewProdMailer(
            os.Getenv("SMTP_FROM"),
            os.Getenv("SMTP_PASSWORD"),
            os.Getenv("SMTP_HOST"),
            os.Getenv("SMTP_PORT"),
        )
    } else {
        cfg.Mailer = mailer.NewDevMailer()
    }

    return cfg
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
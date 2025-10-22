package database

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "os"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        // 環境変数が設定されていない場合のデフォルト値
        dsn = "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
    }
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    DB = db
    return db, nil
}
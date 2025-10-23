// models/password_reset.go
package models

import (
    "time"
)

type PasswordReset struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"not null"`
    Token     string    `gorm:"uniqueIndex;not null"`
    Used      bool      `gorm:"default:false"`
    ExpiresAt time.Time `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
    
    User      User      `gorm:"foreignKey:UserID" json:"-"`
}
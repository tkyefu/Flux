// models/password_reset.go
package models

import (
    "time"
)

type PasswordReset struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    UserID    uint      `gorm:"not null" json:"user_id"`
    Token     string    `gorm:"size:255;not null;uniqueIndex" json:"token"`
    ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
    
    User      User      `gorm:"foreignKey:UserID" json:"-"`
}
package models

import (
    "errors"
    "flux/utils"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

// UserAuth ユーザー認証用のリクエストボディ
type UserAuth struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"` // min=6 から min=8 に変更
}

// ChangePasswordInput パスワード変更用のリクエストボディ
type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}


// HashPassword パスワードをハッシュ化
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword パスワードを検証
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// SetPassword 新しいパスワードを設定
func (u *User) SetPassword(newPassword string) error {
	u.Password = newPassword
	return u.HashPassword()
}

// ChangePassword パスワードを変更
func (u *User) ChangePassword(currentPassword, newPassword string) error {
	// 現在のパスワードを検証
	if err := u.CheckPassword(currentPassword); err != nil {
		return errors.New("現在のパスワードが正しくありません")
	}

	// 新しいパスワードの検証
	if err := utils.ValidatePassword(newPassword, u.Email); err != nil {
		return err
	}

	// 新しいパスワードを設定
	return u.SetPassword(newPassword)
}

// GenerateJWT JWTトークンを生成
func (u *User) GenerateJWT() (string, error) {
    return utils.GenerateToken(u.ID, u.Email)
}

// BeforeCreate ユーザー作成前にパスワードをハッシュ化
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.HashPassword()
}

// AfterCreate ユーザー作成後にパスワードをクリア
func (u *User) AfterCreate(tx *gorm.DB) error {
	u.Password = "" // セキュリティのため、パスワードをクリア
	return nil
}
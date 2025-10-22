package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "flux/models"
    "flux/mailer"
    "flux/utils"
)

type PasswordResetHandler struct {
    DB     *gorm.DB
    Mailer mailer.Mailer
}

func NewPasswordResetHandler(db *gorm.DB, mailer mailer.Mailer) *PasswordResetHandler {
    return &PasswordResetHandler{
        DB:     db,
        Mailer: mailer,
    }
}

func (h *PasswordResetHandler) RequestReset(c *gin.Context) {
    var input struct {
        Email string `json:"email" binding:"required,email"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "有効なメールアドレスを入力してください"})
        return
    }

    // ユーザーを検索
    var user models.User
    if err := h.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
        // セキュリティ上の理由で、ユーザーが存在しなくても成功レスポンスを返す
        c.JSON(http.StatusOK, gin.H{"message": "パスワードリセットの手順を送信しました"})
        return
    }

    // トークン生成と保存
    token := utils.GenerateRandomString(32)
    expiresAt := time.Now().Add(1 * time.Hour)

    resetToken := models.PasswordReset{
        UserID:    user.ID,
        Token:     token,
        ExpiresAt: expiresAt,
    }

    // 古いトークンを削除
    h.DB.Where("user_id = ?", user.ID).Delete(&models.PasswordReset{})

    // 新しいトークンを保存
    if err := h.DB.Create(&resetToken).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "内部エラーが発生しました"})
        return
    }

    // メール送信
    if err := h.Mailer.SendPasswordReset(user.Email, user.Name, token); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "メールの送信に失敗しました"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "パスワードリセットの手順を送信しました"})
}

func (h *PasswordResetHandler) ResetPassword(c *gin.Context) {
    var input struct {
        Token           string `json:"token" binding:"required"`
        NewPassword     string `json:"new_password" binding:"required,min=8"`
        ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
        return
    }

    // トークンを検証
    var resetToken models.PasswordReset
    if err := h.DB.Where("token = ? AND expires_at > ?", input.Token, time.Now()).
        First(&resetToken).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "無効または期限切れのトークンです"})
        return
    }

    // パスワードをハッシュ化
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "パスワードの処理中にエラーが発生しました"})
        return
    }

    // ユーザーのパスワードを更新
    if err := h.DB.Model(&models.User{}).Where("id = ?", resetToken.UserID).
        Update("password", string(hashedPassword)).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "パスワードの更新に失敗しました"})
        return
    }

    // 使用済みトークンを削除
    h.DB.Delete(&resetToken)

    c.JSON(http.StatusOK, gin.H{"message": "パスワードが正常にリセットされました"})
}
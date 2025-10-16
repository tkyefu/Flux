package handlers

import (
	"flux/middleware"
	"flux/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthHandler 認証ハンドラー
type AuthHandler struct {
	DB *gorm.DB
}

// NewAuthHandler 新しいAuthHandlerを作成
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

// RegisterRequest ユーザー登録リクエスト
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest ログインリクエスト
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse 認証レスポンス
type AuthResponse struct {
	Token string     `json:"token"`
	User  models.User `json:"user"`
}

// Register ユーザー登録
func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // メールアドレスの重複チェック
    var existingUser models.User
    if err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "このメールアドレスは既に使用されています"})
        return
    }

    // ユーザー作成
    user := models.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password, // BeforeCreateフックでハッシュ化される
    }

    if err := h.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー登録に失敗しました"})
        return
    }

    // パスワードをクリアしてからレスポンスに含める
    user.Password = ""

    // トークン生成
    token, err := user.GenerateJWT()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "トークンの生成に失敗しました"})
        return
    }

    c.JSON(http.StatusCreated, AuthResponse{
        Token: token,
        User:  user,
    })
}

// Login ログイン
func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // ユーザー検索
    var user models.User
    if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "メールアドレスまたはパスワードが正しくありません"})
        return
    }

    // パスワード検証
    if err := user.CheckPassword(req.Password); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "メールアドレスまたはパスワードが正しくありません"})
        return
    }

    // トークン生成
    token, err := user.GenerateJWT()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "トークンの生成に失敗しました"})
        return
    }

    // シンプルなレスポンスでテスト
    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "token": token,
        "user": gin.H{
            "id":    user.ID,
            "name":  user.Name,
            "email": user.Email,
        },
    })
}

// GetMe ログインユーザー情報取得
func (h *AuthHandler) GetMe(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
        return
    }

    log.Printf("UserID from context: %v", userID)  // デバッグ用

    var user models.User
    if err := h.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
        return
    }

    // パスワードをクリアしてからレスポンスに含める
    user.Password = ""

    c.JSON(http.StatusOK, user)
}

// ChangePassword パスワード変更
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	var input models.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ユーザーを取得
	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
		return
	}

	// パスワードを変更
	if err := user.ChangePassword(input.CurrentPassword, input.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// データベースを更新
	if err := h.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "パスワードの更新に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "パスワードが正常に変更されました"})
}
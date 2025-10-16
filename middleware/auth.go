package middleware

import (
	"flux/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 認証ミドルウェア
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := utils.GetTokenFromRequest(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証トークンが必要です"})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無効なトークンです"})
			c.Abort()
			return
		}

		// ユーザー情報をコンテキストに保存
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}

// GetUserID コンテキストからユーザーIDを取得
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	id, ok := userID.(uint)
	if !ok {
		return 0, false
	}

	return id, true
}

// GetUserEmail コンテキストからユーザーEmailを取得
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("user_email")
	if !exists {
		return "", false
	}

	e, ok := email.(string)
	if !ok {
		return "", false
	}

	return e, true
}
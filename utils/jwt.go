package utils

import (
	"errors"
	"os"
	"time"

	"flux/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	// トークンの有効期間（24時間）
	tokenExpiration = 24 * time.Hour
)

// GenerateToken JWTトークンを生成
func GenerateToken(user *models.User) (string, error) {
	if jwtSecret == nil || len(jwtSecret) == 0 {
		jwtSecret = []byte("default-secret-key") // 開発用のデフォルト値
	}

	claims := &models.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "flux-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken JWTトークンを検証し、クレームを返す
func ParseToken(tokenString string) (*models.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetTokenFromRequest リクエストからJWTトークンを取得
func GetTokenFromRequest(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.Query("token")
	} else if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	return token
}

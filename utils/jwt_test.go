package utils

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/gin-gonic/gin"
)

func TestGenerateAndParseToken(t *testing.T) {
    tkn, err := GenerateToken(1, "user@example.com")
    if err != nil || tkn == "" {
        t.Fatalf("failed to generate token: %v", err)
    }

    claims, err := ParseToken(tkn)
    if err != nil {
        t.Fatalf("failed to parse token: %v", err)
    }
    if claims.UserID != 1 || claims.Email != "user@example.com" {
        t.Fatalf("unexpected claims: %+v", claims)
    }
    if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
        t.Fatalf("token should have future expiration")
    }
}

func TestParseToken_Invalid(t *testing.T) {
    if _, err := ParseToken("invalid.token" ); err == nil {
        t.Fatalf("expected error for invalid token")
    }
}

func TestGetTokenFromRequest(t *testing.T) {
    gin.SetMode(gin.TestMode)

    // Bearer header
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    req, _ := http.NewRequest(http.MethodGet, "/", nil)
    req.Header.Set("Authorization", "Bearer abc123")
    c.Request = req
    if tok := GetTokenFromRequest(c); tok != "abc123" {
        t.Fatalf("expected abc123, got %s", tok)
    }

    // Query param
    w2 := httptest.NewRecorder()
    c2, _ := gin.CreateTestContext(w2)
    req2, _ := http.NewRequest(http.MethodGet, "/?token=qwerty", nil)
    c2.Request = req2
    if tok := GetTokenFromRequest(c2); tok != "qwerty" {
        t.Fatalf("expected qwerty, got %s", tok)
    }

    // Missing
    w3 := httptest.NewRecorder()
    c3, _ := gin.CreateTestContext(w3)
    req3, _ := http.NewRequest(http.MethodGet, "/", nil)
    c3.Request = req3
    if tok := GetTokenFromRequest(c3); tok != "" {
        t.Fatalf("expected empty token, got %s", tok)
    }
}

package middleware

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "flux/utils"
    "github.com/gin-gonic/gin"
)

func TestAuthMiddleware_NoToken(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.Use(AuthMiddleware())
    r.GET("/protected", func(c *gin.Context) { c.Status(http.StatusOK) })

    w := httptest.NewRecorder()
    req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
    r.ServeHTTP(w, req)

    if w.Code != http.StatusUnauthorized {
        t.Fatalf("expected 401, got %d", w.Code)
    }
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.Use(AuthMiddleware())
    r.GET("/protected", func(c *gin.Context) { c.Status(http.StatusOK) })

    w := httptest.NewRecorder()
    req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
    req.Header.Set("Authorization", "Bearer invalid")
    r.ServeHTTP(w, req)

    if w.Code != http.StatusUnauthorized {
        t.Fatalf("expected 401, got %d", w.Code)
    }
}

func TestAuthMiddleware_Success(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.Use(AuthMiddleware())
    r.GET("/protected", func(c *gin.Context) {
        id, _ := c.Get("user_id")
        email, _ := c.Get("user_email")
        c.JSON(http.StatusOK, gin.H{"id": id, "email": email})
    })

    token, err := utils.GenerateToken(42, "user@example.com")
    if err != nil { t.Fatal(err) }

    w := httptest.NewRecorder()
    req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    r.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }

    var res map[string]any
    if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil { t.Fatal(err) }
    if res["id"].(float64) != 42 { t.Fatalf("unexpected id: %v", res["id"]) }
    if res["email"].(string) != "user@example.com" { t.Fatalf("unexpected email: %v", res["email"]) }
}

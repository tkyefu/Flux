package handlers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "flux/models"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type testMailer struct{
    sent int
    lastEmail string
    lastName string
    lastToken string
}

func (m *testMailer) SendPasswordReset(email, username, token string) error {
    m.sent++
    m.lastEmail = email
    m.lastName = username
    m.lastToken = token
    return nil
}

func newTestDB(t *testing.T) *gorm.DB {
    t.Helper()
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil { t.Fatalf("failed to open test db: %v", err) }
    if err := db.AutoMigrate(&models.User{}, &models.PasswordReset{}); err != nil {
        t.Fatalf("failed to migrate: %v", err)
    }
    return db
}

func performJSONRequest(h gin.HandlerFunc, method string, body interface{}) (*httptest.ResponseRecorder, *gin.Context) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    var buf bytes.Buffer
    if body != nil {
        _ = json.NewEncoder(&buf).Encode(body)
    }
    req, _ := http.NewRequest(method, "/", &buf)
    req.Header.Set("Content-Type", "application/json")
    c.Request = req
    return w, c
}

func TestRequestReset_UserNotFound_ReturnsOK(t *testing.T) {
    db := newTestDB(t)
    m := &testMailer{}
    h := NewPasswordResetHandler(db, m)

    w, c := performJSONRequest(h.RequestReset, http.MethodPost, map[string]string{"email":"nobody@example.com"})
    h.RequestReset(c)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }
}

func TestRequestReset_CreatesTokenAndSendsMail(t *testing.T) {
    db := newTestDB(t)
    // create user
    u := models.User{Email: "user@example.com", Name: "User", Password: "hash"}
    if err := db.Create(&u).Error; err != nil { t.Fatal(err) }

    m := &testMailer{}
    h := NewPasswordResetHandler(db, m)

    w, c := performJSONRequest(h.RequestReset, http.MethodPost, map[string]string{"email":u.Email})
    h.RequestReset(c)

    if w.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w.Code) }
    if m.sent != 1 { t.Fatalf("expected mail sent once, got %d", m.sent) }

    var pr models.PasswordReset
    if err := db.Where("user_id = ?", u.ID).First(&pr).Error; err != nil {
        t.Fatalf("expected password reset created: %v", err)
    }
    if pr.Token == "" { t.Fatal("expected token to be set") }
    if pr.ExpiresAt.Before(time.Now()) { t.Fatal("expected future expiration") }
}

func TestResetPassword_Success(t *testing.T) {
    db := newTestDB(t)
    u := models.User{Email: "user@example.com", Name: "User", Password: "oldhash"}
    if err := db.Create(&u).Error; err != nil { t.Fatal(err) }

    token := "tok123"
    pr := models.PasswordReset{UserID: u.ID, Token: token, ExpiresAt: time.Now().Add(time.Hour)}
    if err := db.Create(&pr).Error; err != nil { t.Fatal(err) }

    m := &testMailer{}
    h := NewPasswordResetHandler(db, m)

    body := map[string]string{
        "token": token,
        "new_password": "Password1!",
        "confirm_password": "Password1!",
    }
    w, c := performJSONRequest(h.ResetPassword, http.MethodPost, body)
    h.ResetPassword(c)

    if w.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w.Code) }

    // token used
    var pr2 models.PasswordReset
    if err := db.Where("token = ?", token).First(&pr2).Error; err != nil { t.Fatal(err) }
    if !pr2.Used { t.Fatal("expected token to be marked used") }

    // password updated and hashed
    var u2 models.User
    if err := db.First(&u2, u.ID).Error; err != nil { t.Fatal(err) }
    if bcrypt.CompareHashAndPassword([]byte(u2.Password), []byte("Password1!")) != nil {
        t.Fatal("expected password to match new value")
    }
}

func TestResetPassword_InvalidToken(t *testing.T) {
    db := newTestDB(t)
    u := models.User{Email: "user@example.com", Name: "User", Password: "hash"}
    if err := db.Create(&u).Error; err != nil { t.Fatal(err) }

    m := &testMailer{}
    h := NewPasswordResetHandler(db, m)

    body := map[string]string{
        "token": "nope",
        "new_password": "Password1!",
        "confirm_password": "Password1!",
    }
    w, c := performJSONRequest(h.ResetPassword, http.MethodPost, body)
    h.ResetPassword(c)
    if w.Code != http.StatusBadRequest { t.Fatalf("expected 400, got %d", w.Code) }
}

func TestResetPassword_ExpiredToken(t *testing.T) {
    db := newTestDB(t)
    u := models.User{Email: "user@example.com", Name: "User", Password: "hash"}
    if err := db.Create(&u).Error; err != nil { t.Fatal(err) }

    token := "expired"
    pr := models.PasswordReset{UserID: u.ID, Token: token, ExpiresAt: time.Now().Add(-time.Minute)}
    if err := db.Create(&pr).Error; err != nil { t.Fatal(err) }

    m := &testMailer{}
    h := NewPasswordResetHandler(db, m)

    body := map[string]string{
        "token": token,
        "new_password": "Password1!",
        "confirm_password": "Password1!",
    }
    w, c := performJSONRequest(h.ResetPassword, http.MethodPost, body)
    h.ResetPassword(c)
    if w.Code != http.StatusBadRequest { t.Fatalf("expected 400 for expired token, got %d", w.Code) }
}

func TestResetPassword_UsedToken(t *testing.T) {
    db := newTestDB(t)
    u := models.User{Email: "user@example.com", Name: "User", Password: "hash"}
    if err := db.Create(&u).Error; err != nil { t.Fatal(err) }

    token := "used"
    pr := models.PasswordReset{UserID: u.ID, Token: token, ExpiresAt: time.Now().Add(time.Hour), Used: true}
    if err := db.Create(&pr).Error; err != nil { t.Fatal(err) }

    m := &testMailer{}
    h := NewPasswordResetHandler(db, m)

    body := map[string]string{
        "token": token,
        "new_password": "Password1!",
        "confirm_password": "Password1!",
    }
    w, c := performJSONRequest(h.ResetPassword, http.MethodPost, body)
    h.ResetPassword(c)
    if w.Code != http.StatusBadRequest { t.Fatalf("expected 400 for used token, got %d", w.Code) }
}

func TestResetPassword_WeakPassword(t *testing.T) {
    db := newTestDB(t)
    u := models.User{Email: "user@example.com", Name: "User", Password: "hash"}
    if err := db.Create(&u).Error; err != nil { t.Fatal(err) }

    token := "weak"
    pr := models.PasswordReset{UserID: u.ID, Token: token, ExpiresAt: time.Now().Add(time.Hour)}
    if err := db.Create(&pr).Error; err != nil { t.Fatal(err) }

    m := &testMailer{}
    h := NewPasswordResetHandler(db, m)

    // 9文字の小文字のみ（入力バインディングのmin=8は通るが、複雑性で弾かれる）
    body := map[string]string{
        "token": token,
        "new_password": "simplepwd",
        "confirm_password": "simplepwd",
    }
    w, c := performJSONRequest(h.ResetPassword, http.MethodPost, body)
    h.ResetPassword(c)
    if w.Code != http.StatusBadRequest { t.Fatalf("expected 400 for weak password, got %d", w.Code) }
}

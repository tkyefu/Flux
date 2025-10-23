package handlers

import (
    "encoding/json"
    "net/http"
    "testing"

    "flux/models"
    "golang.org/x/crypto/bcrypt"
)

func TestRegister_Success(t *testing.T) {
    db := newTestDB(t)
    h := NewAuthHandler(db)

    body := RegisterRequest{Name: "User", Email: "user1@example.com", Password: "Password1!"}
    w, c := performJSONRequest(h.Register, http.MethodPost, body)
    h.Register(c)

    if w.Code != http.StatusCreated {
        t.Fatalf("expected 201, got %d", w.Code)
    }

    var res AuthResponse
    if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil { t.Fatal(err) }
    if res.Token == "" { t.Fatal("expected token in response") }
    if res.User.Email != body.Email { t.Fatalf("unexpected email: %s", res.User.Email) }
}

func TestRegister_DuplicateEmail(t *testing.T) {
    db := newTestDB(t)
    if err := db.Create(&models.User{Name: "U", Email: "dup@example.com", Password: "Password1!"}).Error; err != nil { t.Fatal(err) }

    h := NewAuthHandler(db)
    body := RegisterRequest{Name: "U", Email: "dup@example.com", Password: "Password1!"}
    w, c := performJSONRequest(h.Register, http.MethodPost, body)
    h.Register(c)

    if w.Code != http.StatusBadRequest {
        t.Fatalf("expected 400, got %d", w.Code)
    }
}

func TestLogin_SuccessAndFail(t *testing.T) {
    db := newTestDB(t)
    // create user via GORM hook to hash password
    u := models.User{Name: "User", Email: "login@example.com", Password: "Password1!"}
    if err := db.Create(&u).Error; err != nil { t.Fatal(err) }

    h := NewAuthHandler(db)

    // success
    w1, c1 := performJSONRequest(h.Login, http.MethodPost, LoginRequest{Email: u.Email, Password: "Password1!"})
    h.Login(c1)
    if w1.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w1.Code) }

    // fail
    w2, c2 := performJSONRequest(h.Login, http.MethodPost, LoginRequest{Email: u.Email, Password: "wrong"})
    h.Login(c2)
    if w2.Code != http.StatusUnauthorized { t.Fatalf("expected 401, got %d", w2.Code) }
}

func TestGetMe_Success(t *testing.T) {
    db := newTestDB(t)
    u := models.User{Name: "Me", Email: "me@example.com", Password: "Password1!"}
    if err := db.Create(&u).Error; err != nil { t.Fatal(err) }

    h := NewAuthHandler(db)
    w, c := performJSONRequest(h.GetMe, http.MethodGet, nil)
    c.Set("user_id", u.ID)

    h.GetMe(c)
    if w.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w.Code) }
}

func TestChangePassword_Success(t *testing.T) {
    db := newTestDB(t)
    u := models.User{Name: "Ch", Email: "ch@example.com", Password: "Password1!"}
    if err := db.Create(&u).Error; err != nil { t.Fatal(err) }

    h := NewAuthHandler(db)
    body := models.ChangePasswordInput{CurrentPassword: "Password1!", NewPassword: "NewPass1!"}
    w, c := performJSONRequest(h.ChangePassword, http.MethodPost, body)
    c.Set("user_id", u.ID)
    h.ChangePassword(c)

    if w.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w.Code) }

    var u2 models.User
    if err := db.First(&u2, u.ID).Error; err != nil { t.Fatal(err) }
    if bcrypt.CompareHashAndPassword([]byte(u2.Password), []byte("NewPass1!")) != nil {
        t.Fatal("expected password to be updated")
    }
}

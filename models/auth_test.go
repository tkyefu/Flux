package models

import (
    "testing"

    "flux/utils"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func TestHashAndCheckPassword(t *testing.T) {
    u := &User{Password: "Password1!"}
    if err := u.HashPassword(); err != nil { t.Fatal(err) }
    if err := u.CheckPassword("Password1!"); err != nil { t.Fatalf("expected password to match: %v", err) }
    if err := u.CheckPassword("wrong"); err == nil { t.Fatalf("expected mismatch error") }
}

func TestSetAndChangePassword(t *testing.T) {
    u := &User{Password: "OldPass1!"}
    if err := u.HashPassword(); err != nil { t.Fatal(err) }

    // change wrong current
    if err := u.ChangePassword("bad", "NewPass1!"); err == nil { t.Fatalf("expected error for wrong current") }

    // weak new
    if err := u.ChangePassword("OldPass1!", "simplepwd"); err == nil { t.Fatalf("expected error for weak password") }

    // success
    if err := u.ChangePassword("OldPass1!", "NewPass1!"); err != nil { t.Fatalf("change password failed: %v", err) }
    if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("NewPass1!")) != nil { t.Fatalf("password not updated") }
}

func TestGormHooksAndJWT(t *testing.T) {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil { t.Fatal(err) }
    if err := db.AutoMigrate(&User{}); err != nil { t.Fatal(err) }

    u := User{Name: "U", Email: "u@example.com", Password: "Password1!"}
    if err := db.Create(&u).Error; err != nil { t.Fatal(err) }

    // AfterCreate should clear password in the struct
    if u.Password != "" { t.Fatalf("expected password cleared after create, got %q", u.Password) }

    var stored User
    if err := db.First(&stored, u.ID).Error; err != nil { t.Fatal(err) }
    if err := stored.CheckPassword("Password1!"); err != nil { t.Fatalf("stored password should be hashed: %v", err) }

    // JWT
    token, err := u.GenerateJWT()
    if err != nil || token == "" { t.Fatalf("failed to generate jwt: %v", err) }
    claims, err := utils.ParseToken(token)
    if err != nil { t.Fatalf("failed to parse jwt: %v", err) }
    if claims.UserID != u.ID || claims.Email != u.Email { t.Fatalf("unexpected claims: %+v", claims) }
}

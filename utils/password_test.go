package utils

import (
    "errors"
    "testing"
)

func resetPasswordEnv(t *testing.T) {
    t.Helper()
    t.Setenv("PASSWORD_MIN_LENGTH", "")
    t.Setenv("PASSWORD_REQUIRE_UPPER", "")
    t.Setenv("PASSWORD_REQUIRE_LOWER", "")
    t.Setenv("PASSWORD_REQUIRE_NUMBER", "")
    t.Setenv("PASSWORD_REQUIRE_SPECIAL", "")
}

func TestValidatePassword_DefaultPolicy(t *testing.T) {
    resetPasswordEnv(t)

    if err := ValidatePassword("short", "user@example.com"); !errors.Is(err, ErrPasswordTooShort) {
        t.Fatalf("expected ErrPasswordTooShort, got %v", err)
    }

    if err := ValidatePassword("onlylower", "user@example.com"); !errors.Is(err, ErrPasswordWeak) {
        t.Fatalf("expected ErrPasswordWeak, got %v", err)
    }

    if err := ValidatePassword("Password1!", "user@example.com"); err != nil {
        t.Fatalf("expected password to be valid, got %v", err)
    }
}

func TestValidatePassword_CustomPolicy(t *testing.T) {
    resetPasswordEnv(t)
    t.Setenv("PASSWORD_MIN_LENGTH", "12")
    t.Setenv("PASSWORD_REQUIRE_UPPER", "true")
    t.Setenv("PASSWORD_REQUIRE_LOWER", "true")
    t.Setenv("PASSWORD_REQUIRE_NUMBER", "true")
    t.Setenv("PASSWORD_REQUIRE_SPECIAL", "true")

    if err := ValidatePassword("short1!A", "user@example.com"); !errors.Is(err, ErrPasswordTooShort) {
        t.Fatalf("expected ErrPasswordTooShort, got %v", err)
    }

    if err := ValidatePassword("longpassword1", "user@example.com"); !errors.Is(err, ErrPasswordWeak) {
        t.Fatalf("expected ErrPasswordWeak when missing required classes, got %v", err)
    }

    if err := ValidatePassword("StrongPass1!", "user@example.com"); err != nil {
        t.Fatalf("expected password to satisfy custom policy, got %v", err)
    }
}

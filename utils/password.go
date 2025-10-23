package utils

import (
	"os"
	"strconv"
	"errors"
	"strings"
	"unicode"
)

var (
	ErrPasswordTooShort      = errors.New("パスワードは8文字以上である必要があります")
	ErrPasswordTooCommon     = errors.New("このパスワードは使用できません。より複雑なパスワードを設定してください")
	ErrPasswordContainsEmail = errors.New("パスワードにメールアドレスを含めることはできません")
	ErrPasswordWeak          = errors.New("パスワードは大文字・小文字・数字・記号のうち3種類以上を含む必要があります")
)

// よく使われるパスワードのリスト
var commonPasswords = map[string]bool{
	"password":     true,
	"12345678":     true,
	"qwertyui":     true,
	"admin123":     true,
	"welcome123":   true,
	"password123":  true,
	"123456789":    true,
	"qwerty123":    true,
	"admin1234":    true,
	"welcome1":     true,
}

// ValidatePassword はパスワードがポリシーを満たしているか検証します
func ValidatePassword(password, email string) error {
	minLen := getEnvInt("PASSWORD_MIN_LENGTH", 8)
	if len(password) < minLen {
		return ErrPasswordTooShort
	}

	if isCommonPassword(password) {
		return ErrPasswordTooCommon
	}

	if email != "" && strings.Contains(strings.ToLower(password), strings.Split(email, "@")[0]) {
		return ErrPasswordContainsEmail
	}

	reqUpper := getEnvBool("PASSWORD_REQUIRE_UPPER", false)
	reqLower := getEnvBool("PASSWORD_REQUIRE_LOWER", false)
	reqNumber := getEnvBool("PASSWORD_REQUIRE_NUMBER", false)
	reqSpecial := getEnvBool("PASSWORD_REQUIRE_SPECIAL", false)

	if reqUpper || reqLower || reqNumber || reqSpecial {
		hasUpper, hasLower, hasNumber, hasSpecial := classify(password)
		if reqUpper && !hasUpper { return ErrPasswordWeak }
		if reqLower && !hasLower { return ErrPasswordWeak }
		if reqNumber && !hasNumber { return ErrPasswordWeak }
		if reqSpecial && !hasSpecial { return ErrPasswordWeak }
	} else if !meetsComplexityRequirements(password) {
		return ErrPasswordWeak
	}

	return nil
}

// isCommonPassword はパスワードが一般的なパスワードでないかチェックします
func isCommonPassword(password string) bool {
	return commonPasswords[strings.ToLower(password)]
}

// meetsComplexityRequirements はパスワードの複雑性要件を満たしているかチェックします
func meetsComplexityRequirements(password string) bool {
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}

		// 3つ以上のカテゴリに含まれていれば早期リターン
		categories := 0
		if hasUpper {
			categories++
		}
		if hasLower {
			categories++
		}
		if hasNumber {
			categories++
		}
		if hasSpecial {
			categories++
		}
		if categories >= 3 {
			return true
		}
	}

	return false
}

func classify(password string) (bool, bool, bool, bool) {
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasUpper, hasLower, hasNumber, hasSpecial
}

func getEnvInt(name string, def int) int {
	if v := os.Getenv(name); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i > 0 {
			return i
		}
	}
	return def
}

func getEnvBool(name string, def bool) bool {
	if v := os.Getenv(name); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}
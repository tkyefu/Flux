package utils

import (
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
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	if isCommonPassword(password) {
		return ErrPasswordTooCommon
	}

	if email != "" && strings.Contains(strings.ToLower(password), strings.Split(email, "@")[0]) {
		return ErrPasswordContainsEmail
	}

	if !meetsComplexityRequirements(password) {
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
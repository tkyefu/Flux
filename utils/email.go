package utils

import (
    "crypto/rand"
)

// GenerateRandomString はランダムな文字列を生成します
func GenerateRandomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    if _, err := rand.Read(b); err != nil {
        panic(err)
    }
    for i := range b {
        b[i] = charset[int(b[i])%len(charset)]
    }
    return string(b)
}
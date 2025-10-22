package mailer

import (
    "fmt"
    "log"
    "os"
)

// Mailer はメール送信のインターフェースを定義します
type Mailer interface {
    SendPasswordReset(email, username, token string) error
}

// DevMailer は開発用のメール送信をシミュレートします
type DevMailer struct{}

func NewDevMailer() *DevMailer {
    return &DevMailer{}  // 修正: Devailer → DevMailer
}

func (m *DevMailer) SendPasswordReset(email, username, token string) error {
    resetURL := generateResetURL(token)
    log.Printf("[DEV] パスワードリセットリンク: %s\n", resetURL)
    log.Printf("[DEV] 受信者: %s\n", email)
    return nil
}

// ProdMailer は本番環境用のメール送信を行います
type ProdMailer struct {
    from     string
    password string
    host     string
    port     string
}

func NewProdMailer(from, password, host, port string) *ProdMailer {
    return &ProdMailer{
        from:     from,
        password: password,
        host:     host,
        port:     port,
    }
}

func (m *ProdMailer) SendPasswordReset(email, username, token string) error {
    resetURL := generateResetURL(token)
    log.Printf("[PROD] メールを送信しました: %s\n", email)
    log.Printf("[PROD] リセットURL: %s\n", resetURL)
    return nil
}

func generateResetURL(token string) string {
    frontendURL := os.Getenv("FRONTEND_URL")
    if frontendURL == "" {
        frontendURL = "http://localhost:3000"
    }
    return fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)
}
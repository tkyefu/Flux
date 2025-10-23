#!/bin/bash

# テスト用のメールアドレス
TEST_EMAIL="test@example.com"

# パスワードリセットのリクエスト
echo "Requesting password reset..."
curl -X POST http://localhost:8080/api/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$TEST_EMAIL\"}"

# データベースからトークンを取得
echo -e "\n\nGetting token from database..."
TOKEN_QUERY="SELECT token FROM password_resets ORDER BY created_at DESC LIMIT 1;"
TOKEN=$(psql -U flux_user -d flux -t -c "$TOKEN_QUERY" | tr -d '[:space:]')

if [ -z "$TOKEN" ]; then
    echo -e "\nError: No token found in database"
    exit 1
fi

echo -e "\nResetting password with token: $TOKEN"
curl -X POST http://localhost:8080/api/auth/reset-password \
  -H "Content-Type: application/json" \
  -d "{\"token\": \"$TOKEN\", \"new_password\": \"NewSecurePassword123!\", \"confirm_password\": \"NewSecurePassword123!\"}"

echo -e "\n\nTrying to use the same token again..."
curl -X POST http://localhost:8080/api/auth/reset-password \
  -H "Content-Type: application/json" \
  -d "{\"token\": \"$TOKEN\", \"new_password\": \"AnotherPassword123!\", \"confirm_password\": \"AnotherPassword123!\"}"
#!/bin/bash

# エラーが発生したら即時終了
set -e

# 環境変数の読み込み
if [ -f .env ]; then
    # 空行と#で始まる行を除外してからexport
    while IFS= read -r line; do
        # 空行または#で始まる行はスキップ
        if [[ -n "$line" && ! "$line" =~ ^[[:space:]]*# ]]; then
            export "$line"
        fi
    done < .env
fi

# デフォルト値の設定
DB_USER=${DB_USER:-flux_user}
DB_NAME=${DB_NAME:-flux}
DB_HOST=${DB_HOST:-}
DB_PORT=${DB_PORT:-}
API_URL=${API_URL:-http://localhost:8080}

# psql 接続オプション（任意のホスト/ポートに対応）
PSQLOPTS=""
if [ -n "$DB_HOST" ]; then PSQLOPTS+=" -h $DB_HOST"; fi
if [ -n "$DB_PORT" ]; then PSQLOPTS+=" -p $DB_PORT"; fi

# テスト用のユーザー情報（環境変数で上書き可能）
TEST_EMAIL=${TEST_EMAIL:-"test@example.com"}
TEST_PASSWORD=${TEST_PASSWORD:-"NewSecurePassword123!"}

# ヘルパー関数
log_success() {
    echo -e "\033[0;32m[SUCCESS] $1\033[0m"
}

log_error() {
    echo -e "\033[0;31m[ERROR] $1\033[0m" >&2
    exit 1
}

# テストデータのセットアップ
setup_test_data() {
    echo "Setting up test data..."
    # 既存のテストデータをクリーンアップ
    psql $PSQLOPTS -U $DB_USER -d $DB_NAME -c "DELETE FROM password_resets; DELETE FROM users WHERE email = '$TEST_EMAIL';" > /dev/null 2>&1 || true
    
    # テストユーザーを作成
    psql $PSQLOPTS -U $DB_USER -d $DB_NAME -c "INSERT INTO users (email, name, password) VALUES ('$TEST_EMAIL', 'Test User', 'dummy_hash') ON CONFLICT (email) DO NOTHING;"
}

# リクエスト送信ヘルパー
send_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local expected_status=${4:-200}
    
    echo -e "\nSending $method request to $endpoint"
    if [ -n "$data" ]; then
        response=$(curl -s -o response.txt -w "%{http_code}" -X $method \
            -H "Content-Type: application/json" \
            -d "$data" \
            $API_URL$endpoint)
    else
        response=$(curl -s -o response.txt -w "%{http_code}" -X $method \
            -H "Content-Type: application/json" \
            $API_URL$endpoint)
    fi
    
    if [ "$response" -ne "$expected_status" ]; then
        log_error "Expected status $expected_status but got $response"
        cat response.txt
        return 1
    fi
    
    cat response.txt
    return 0
}

# ステータスコードのみ取得したい場合のヘルパー
send_request_status() {
    local method=$1
    local endpoint=$2
    local data=$3
    if [ -n "$data" ]; then
        curl -s -o /dev/null -w "%{http_code}" -X $method \
            -H "Content-Type: application/json" \
            -d "$data" \
            $API_URL$endpoint
    else
        curl -s -o /dev/null -w "%{http_code}" -X $method \
            -H "Content-Type: application/json" \
            $API_URL$endpoint
    fi
}

# メインのテスト実行
run_tests() {
    echo "Starting password reset tests..."
    
    # 1. パスワードリセットリクエストのテスト
    echo -e "\n--- Testing password reset request ---"
    send_request "POST" "/api/v1/auth/forgot-password" "{\"email\": \"$TEST_EMAIL\"}" 200
    
    # 2. トークンの取得
    echo -e "\n--- Getting reset token ---"
    TOKEN_QUERY="SELECT token FROM password_resets WHERE user_id = (SELECT id FROM users WHERE email = '$TEST_EMAIL') ORDER BY created_at DESC LIMIT 1;"
    TOKEN=$(psql $PSQLOPTS -U $DB_USER -d $DB_NAME -t -c "$TOKEN_QUERY" | tr -d '[:space:]')
    
    if [ -z "$TOKEN" ]; then
        log_error "No token found in database"
    fi
    
    echo "Retrieved token: $TOKEN"
    FIRST_TOKEN="$TOKEN"
    
    # 3. 有効なトークンでのパスワードリセット
    echo -e "\n--- Testing valid password reset ---"
    send_request "POST" "/api/v1/auth/reset-password" \
        "{\"token\": \"$TOKEN\", \"new_password\": \"$TEST_PASSWORD\", \"confirm_password\": \"$TEST_PASSWORD\"}" 200
    
    # 4. パスワードの複雑性チェック（先に実施してレート制限を回避）
    echo -e "\n--- Testing password complexity (should fail) ---"
    send_request "POST" "/api/v1/auth/forgot-password" "{\"email\": \"$TEST_EMAIL\"}" 200
    # 新しいトークンを取得（複雑性テスト用）
    TOKEN=$(psql $PSQLOPTS -U $DB_USER -d $DB_NAME -t -c "$TOKEN_QUERY" | tr -d '[:space:]')
    send_request "POST" "/api/v1/auth/reset-password" \
        "{\"token\": \"$TOKEN\", \"new_password\": \"simple\", \"confirm_password\": \"simple\"}" 400

    # 5. 同じトークンでの再試行（失敗するはず） - 最初のトークンを再利用
    echo -e "\n--- Testing reused token (should fail) ---"
    send_request "POST" "/api/v1/auth/reset-password" \
        "{\"token\": \"$FIRST_TOKEN\", \"new_password\": \"NewPassword123!\", \"confirm_password\": \"NewPassword123!\"}" 400

    # 6. 無効なトークンでのテスト
    echo -e "\n--- Testing invalid token (should fail) ---"
    send_request "POST" "/api/v1/auth/reset-password" \
        "{\"token\": \"invalid-token-123\", \"new_password\": \"NewPassword123!\", \"confirm_password\": \"NewPassword123!\"}" 400
    
    # 7. 期限切れトークンでのテスト
    echo -e "\n--- Testing expired token (should fail) ---"
    EXPIRED_TOKEN=$(openssl rand -hex 16)
    psql $PSQLOPTS -U $DB_USER -d $DB_NAME -c "INSERT INTO password_resets (user_id, token, used, expires_at, created_at, updated_at) VALUES ((SELECT id FROM users WHERE email = '$TEST_EMAIL'), '$EXPIRED_TOKEN', false, NOW() - INTERVAL '1 minute', NOW(), NOW());"
    send_request "POST" "/api/v1/auth/reset-password" \
        "{\"token\": \"$EXPIRED_TOKEN\", \"new_password\": \"NewPassword123!\", \"confirm_password\": \"NewPassword123!\"}" 400
    
    # 8. レートリミットのテスト（環境変数に応じて動的に検証）
    echo -e "\n--- Testing rate limiting (dynamic) ---"
    LIMIT=${RATE_LIMIT_REQUESTS:-5}
    ATTEMPTS=$((LIMIT + 1))
    rl_429=0
    echo "Configured RATE_LIMIT_REQUESTS=$LIMIT"
    for i in $(seq 1 $ATTEMPTS); do
        code=$(send_request_status "POST" "/api/v1/auth/forgot-password" "{\"email\": \"$TEST_EMAIL\"}")
        echo "Request $i/$ATTEMPTS -> HTTP $code"
        if [ "$code" -eq 429 ]; then rl_429=1; fi
    done
    if [ "$LIMIT" -le 10 ]; then
        if [ $rl_429 -ne 1 ]; then
            log_error "Rate limiting did not trigger as expected (limit=$LIMIT)"
        fi
    else
        if [ $rl_429 -eq 1 ]; then
            log_success "Rate limiting triggered even with high limit ($LIMIT)"
        else
            echo "Skipping strict rate limit assertion (limit too high: $LIMIT)"
        fi
    fi
    
    log_success "All tests passed successfully!"
}

# メイン実行
setup_test_data
run_tests

# 後片付け（テストデータの削除）
teardown_test_data() {
    echo "Tearing down test data..."
    psql $PSQLOPTS -U $DB_USER -d $DB_NAME -c "DELETE FROM password_resets WHERE user_id = (SELECT id FROM users WHERE email = '$TEST_EMAIL'); DELETE FROM users WHERE email = '$TEST_EMAIL';" > /dev/null 2>&1 || true
}
teardown_test_data
rm -f response.txt
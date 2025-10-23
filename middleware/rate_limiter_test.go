package middleware

import (
    "testing"
    "time"
)

func TestGetVisitor_Defaults(t *testing.T) {
    t.Setenv("RATE_LIMIT_REQUESTS", "")
    t.Setenv("RATE_LIMIT_WINDOW", "")

    limiter := getVisitor("127.0.0.1")
    if limiter == nil {
        t.Fatal("expected limiter, got nil")
    }

    // 5 requests should be allowed in default window
    for i := 0; i < 5; i++ {
        if !limiter.Allow() {
            t.Fatalf("request %d should be allowed with default settings", i+1)
        }
    }

    if limiter.Allow() {
        t.Fatal("6th request should not be allowed with default settings")
    }
}

func TestGetVisitor_CustomConfig(t *testing.T) {
    t.Setenv("RATE_LIMIT_REQUESTS", "2")
    t.Setenv("RATE_LIMIT_WINDOW", "1s")

    limiter := getVisitor("192.168.0.1")
    if limiter == nil {
        t.Fatal("expected limiter, got nil")
    }

    if !limiter.Allow() || !limiter.Allow() {
        t.Fatal("first two requests should be allowed")
    }

    if limiter.Allow() {
        t.Fatal("third request should be blocked")
    }

    time.Sleep(time.Second)

    if !limiter.Allow() {
        t.Fatal("request should be allowed after window reset")
    }
}

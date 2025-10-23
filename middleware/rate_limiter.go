package middleware

import (
    "os"
    "strconv"
    "net/http"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
)

type visitor struct {
    limiter  *rate.Limiter
    lastSeen time.Time
}

var visitors = make(map[string]*visitor)
var mu sync.Mutex

func init() {
    go cleanupVisitors()
}

func getEnvInt(name string, def int) int {
    if v := os.Getenv(name); v != "" {
        if i, err := strconv.Atoi(v); err == nil && i > 0 {
            return i
        }
    }
    return def
}

func getEnvDuration(name string, def time.Duration) time.Duration {
    if v := os.Getenv(name); v != "" {
        if d, err := time.ParseDuration(v); err == nil && d > 0 {
            return d
        }
    }
    return def
}

func getVisitor(ip string) *rate.Limiter {
    mu.Lock()
    defer mu.Unlock()

    v, exists := visitors[ip]
    if !exists {
        requests := getEnvInt("RATE_LIMIT_REQUESTS", 5)
        window := getEnvDuration("RATE_LIMIT_WINDOW", time.Minute)
        per := window / time.Duration(requests)
        if per <= 0 {
            per = time.Minute / 5
        }
        limiter := rate.NewLimiter(rate.Every(per), requests)
        visitors[ip] = &visitor{limiter, time.Now()}
        return limiter
    }
    v.lastSeen = time.Now()
    return v.limiter
}

func cleanupVisitors() {
    for {
        time.Sleep(time.Minute)
        mu.Lock()
        for ip, v := range visitors {
            if time.Since(v.lastSeen) > 3*time.Minute {
                delete(visitors, ip)
            }
        }
        mu.Unlock()
    }
}

func RateLimit() gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()
        limiter := getVisitor(ip)
        
        if !limiter.Allow() {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "リクエストが多すぎます。しばらく待ってからお試しください。",
            })
            c.Abort()
            return
        }
        c.Next()
    }
}
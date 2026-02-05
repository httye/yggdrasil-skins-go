package middleware

import (
	"sync"
	"time"

	"github.com/httye/yggdrasil-skins-go/src/utils"

	"github.com/gin-gonic/gin"
)

// RateLimiter 速率限制�?type RateLimiter struct {
	requests map[string]time.Time // 记录每个用户的最后请求时�?	mu       sync.RWMutex         // 读写�?	interval time.Duration        // 请求间隔
}

// NewRateLimiter 创建新的速率限制�?func NewRateLimiter(interval time.Duration) *RateLimiter {
	limiter := &RateLimiter{
		requests: make(map[string]time.Time),
		interval: interval,
	}

	// 启动清理协程
	go limiter.cleanup()
	return limiter
}

// cleanup 定期清理过期的请求记�?func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, lastTime := range rl.requests {
			if now.Sub(lastTime) > rl.interval*2 {
				delete(rl.requests, key)
			}
		}
		rl.mu.Unlock()
	}
}

// Allow 检查是否允许请�?func (rl *RateLimiter) Allow(identifier string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	lastTime, exists := rl.requests[identifier]

	if exists && now.Sub(lastTime) < rl.interval {
		return false
	}

	rl.requests[identifier] = now
	return true
}

// RateLimit 速率限制中间件（性能优化版）
func RateLimit(interval time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(interval)

	return func(c *gin.Context) {
		// 使用客户端IP作为标识符（避免消耗请求体�?		identifier := c.ClientIP()

		if !limiter.Allow(identifier) {
			// 使用缓存的错误响应，避免重复序列�?			utils.RespondCachedError(c, 403, "rate_limit_exceeded")
			c.Abort()
			return
		}

		c.Next()
	}
}

package middleware

import (
	"context"
	"time"

	"clean-arch-go/internal/pkg/redis"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter handles rate limiting
func NewRateLimiter(redisClient *redis.RedisClient, limit rate.Limit, burst int) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		limiter:     rate.NewLimiter(limit, burst),
	}
}

type RateLimiter struct {
	redisClient *redis.RedisClient
	limiter     *rate.Limiter
}

// Limit is a middleware that limits request rate
func (r *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		key := "rate_limit:" + clientIP

		// Check Redis for existing rate limit
		if err := r.redisClient.Set(context.Background(), key, "1", 10*time.Second); err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		// Check rate limiter
		if !r.limiter.Allow() {
			c.JSON(429, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}

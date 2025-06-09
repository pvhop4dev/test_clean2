package middleware

import (
	"clean-arch-go/internal/pkg/redis"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// NewLogger returns a new logger middleware
func NewLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		if len(c.Errors) > 0 {
			c.Writer.Header().Set("X-Error", "true")
		}

		c.Writer.Header().Set("X-Response-Time", latency.String())
	}
}

// NewRecovery returns a new recovery middleware
func NewRecovery() gin.HandlerFunc {
	return gin.Recovery()
}

// NewCORS returns a new CORS middleware
func NewCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

// NewRateLimiter returns a new rate limiter middleware
func NewRateLimiter(redisClient *redis.RedisClient, limit int, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement simple rate limiting using Redis or in-memory
		// For now, just pass through
		c.Next()
	}
}

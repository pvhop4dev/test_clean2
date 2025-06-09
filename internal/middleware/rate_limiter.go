package middleware

import (
	"context"
	"net/http"
	"time"

	"clean-arch-go/internal/pkg/redis"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	redisClient *redis.RedisClient
	limiter     *rate.Limiter
}

func NewRateLimiter(redisClient *redis.RedisClient, limit rate.Limit, burst int) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		limiter:     rate.NewLimiter(limit, burst),
	}
}

func (r *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sử dụng Redis để lưu trữ rate limit dựa trên IP
		ip := c.ClientIP()
		key := "rate_limit:" + ip

		// Tăng số lượng request
		ctx := context.Background()
		count, err := r.redisClient.Incr(ctx, key)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		// Đặt thời gian hết hạn cho key nếu đây là lần đầu tiên
		if count == 1 {
			r.redisClient.Expire(ctx, key, time.Minute)
		}

		// Kiểm tra số lượng request vượt quá giới hạn
		if count > 100 { // Giới hạn 100 request mỗi phút
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "too many requests",
				"retry_after": 60, // giây
			})
			return
		}

		// Sử dụng rate limiter của thư viện
		if !r.limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "too many requests",
				"retry_after": 1, // giây
			})
			return
		}

		c.Next()
	}
}

// RedisClient mở rộng để hỗ trợ rate limiting
type RedisClient interface {
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
}

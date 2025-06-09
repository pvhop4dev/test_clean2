package middleware

import (
	"clean-arch-go/internal/domain/service"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware handles authentication
func NewAuthMiddleware(authSvc service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authSvc: authSvc,
	}
}

type AuthMiddleware struct {
	authSvc service.AuthService
}

// AuthRequired is a middleware that checks for authentication
func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, err := m.authSvc.ValidateToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims)
		c.Next()
	}
}

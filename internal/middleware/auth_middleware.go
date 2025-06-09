package middleware

import (
	"net/http"
	"strings"

	"clean-arch-go/internal/delivery/http/helper"
	"clean-arch-go/internal/domain/service"
	"clean-arch-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
)

type AuthMiddleware struct {
	authService service.AuthService
}

func NewAuthMiddleware(authService service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			msg := helper.MustTranslate(c, "error.unauthorized", nil)
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			msg := helper.MustTranslate(c, "error.invalid_token_format", nil)
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			c.Abort()
			return
		}

		userID, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			logger.Error("Failed to validate token: %v", err)
			msg := helper.MustTranslate(c, "error.invalid_token", nil)
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			c.Abort()
			return
		}

		// Set user ID in context for handlers to use
		c.Set("userID", userID)
		c.Next()
	}
}

package middleware

import (
	"clean-arch-go/internal/domain/service"
	"clean-arch-go/internal/pkg/i18n"
	"clean-arch-go/internal/pkg/logger"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix       = "Bearer "
)

type AuthMiddleware struct {
	authService service.AuthService
	localizer   *i18n.Localizer
}

func NewAuthMiddleware(authService service.AuthService, localizer *i18n.Localizer) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		localizer:   localizer,
	}
}

func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errMsg, _ := m.localizer.Translate(c.Request.Context(), "auth.unauthorized", nil)
			c.JSON(http.StatusUnauthorized, gin.H{"error": errMsg})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			errMsg, _ := m.localizer.Translate(c.Request.Context(), "auth.invalid_token_format", nil)
			c.JSON(http.StatusUnauthorized, gin.H{"error": errMsg})
			c.Abort()
			return
		}

		userID, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			logger.Error("Failed to validate token: %v", err)
			errMsg, _ := m.localizer.Translate(c.Request.Context(), "auth.invalid_token", nil)
			c.JSON(http.StatusUnauthorized, gin.H{"error": errMsg})
			c.Abort()
			return
		}

		// Set user ID in context for handlers to use
		c.Set("userID", userID)
		c.Next()
	}
}

package middleware

import (
	"clean-arch-go/internal/domain/service"
	"clean-arch-go/internal/domain/user"
	"clean-arch-go/internal/errors"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// UserKey is the key used to store the user in the context
	UserKey = "user"
	// TokenKey is the key used to store the token in the context
	TokenKey = "token"
)

type AuthMiddleware struct {
	authSvc service.AuthService
}

// NewAuthMiddleware creates a new AuthMiddleware instance
func NewAuthMiddleware(authSvc service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authSvc: authSvc,
	}
}

// AuthRequired is a middleware that checks for a valid JWT token in the Authorization header
func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := extractToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Validate token and get user
		user, err := m.authSvc.GetUserFromToken(c.Request.Context(), tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user and token in context
		ctx := context.WithValue(c.Request.Context(), UserKey, user)
		ctx = context.WithValue(ctx, TokenKey, tokenString)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// AuthOptional is a middleware that checks for a valid JWT token but doesn't require it
func (m *AuthMiddleware) AuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := extractToken(c)
		if err != nil || tokenString == "" {
			c.Next()
			return
		}

		// Validate token and get user
		user, err := m.authSvc.GetUserFromToken(c.Request.Context(), tokenString)
		if err == nil && user != nil {
			// Set user and token in context if valid
			ctx := context.WithValue(c.Request.Context(), UserKey, user)
			ctx = context.WithValue(ctx, TokenKey, tokenString)
			c.Request = c.Request.WithContext(ctx)
		}

		c.Next()
	}
}

// extractToken extracts the JWT token from the Authorization header
func extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.NewAppError("UNAUTHORIZED", "Authorization header is required", nil)
	}

	// Format: Bearer <token>
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.NewAppError("UNAUTHORIZED", "Invalid Authorization header format", nil)
	}

	return parts[1], nil
}

// GetUserFromContext returns the authenticated user from the context
func GetUserFromContext(ctx context.Context) (*user.User, bool) {
	user, ok := ctx.Value(UserKey).(*user.User)
	return user, ok
}

// GetTokenFromContext returns the JWT token from the context
func GetTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(TokenKey).(string)
	return token, ok
}

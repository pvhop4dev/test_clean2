package middleware

import (
	"clean-arch-go/internal/domain/service"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	service service.AuthService
}

func NewAuthMiddleware(service service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		service: service,
	}
}

func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement authentication middleware
		c.Next()
	}
}

func (m *AuthMiddleware) AuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement optional authentication middleware
		c.Next()
	}
}

package handler

import (
	"fmt"
	"net/http"
	"time"

	"clean-arch-go/internal/domain/service"
	"clean-arch-go/internal/domain/translation"
	"clean-arch-go/internal/pkg/i18n"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.GetLocalizer().MustTranslate(language.English, translation.ErrInvalidCredentials, nil)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": i18n.GetLocalizer().MustTranslate(language.English, translation.AuthRegisterSuccess, nil),
		"user": AuthResponse{
			ID:        fmt.Sprintf("%d", user.ID),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": i18n.GetLocalizer().MustTranslate(language.English, translation.ErrInvalidCredentials, nil)})
		return
	}

	// Lấy thông tin user từ token
	userID, err := h.authService.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.GetLocalizer().MustTranslate(language.English, translation.ErrInternalServerError, nil)})
		return
	}

	// Lấy thông tin user
	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.GetLocalizer().MustTranslate(language.English, translation.ErrInternalServerError, nil)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": i18n.GetLocalizer().MustTranslate(language.English, translation.AuthLoginSuccess, nil),
		"token":   token,
		"user": AuthResponse{
			ID:        fmt.Sprintf("%d", user.ID),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	})
}

package handler

import (
	"context"
	"net/http"

	"clean-arch-go/internal/domain/service"
	"github.com/gin-gonic/gin"
)

// LoginInput represents the login request body
// swagger:parameters login
type LoginInput struct {
	// Email of the user
	// required: true
	// example: user@example.com
	Email string `json:"email" binding:"required,email"`

	// Password of the user
	// required: true
	// minLength: 8
	// example: password123
	Password string `json:"password" binding:"required,min=8"`
}

// RegisterInput represents the registration request body
// swagger:parameters register
type RegisterInput struct {
	// Name of the user
	// required: true
	// minLength: 2
	// example: John Doe
	Name string `json:"name" binding:"required,min=2"`

	// Email of the user
	// required: true
	// example: user@example.com
	Email string `json:"email" binding:"required,email"`

	// Password of the user
	// required: true
	// minLength: 8
	// example: password123
	Password string `json:"password" binding:"required,min=8"`
}

// TokenResponse represents the authentication token response
// swagger:response tokenResponse
type TokenResponse struct {
	// JWT access token
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	AccessToken string `json:"access_token"`

	// Type of token
	// example: bearer
	TokenType string `json:"token_type"`

	// Expiration time in seconds
	// example: 3600
	ExpiresIn int64 `json:"expires_in"`
}

// ErrorResponse represents an error response
// swagger:response errorResponse
type ErrorResponse struct {
	// Error message
	// example: Invalid credentials
	Error string `json:"error"`
}

type AuthHandler struct {
	authSvc service.AuthService
}

func NewAuthHandler(authSvc service.AuthService) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
	}
}

// RegisterAuthRoutes registers the authentication routes
func (h *AuthHandler) RegisterAuthRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
	}
}

// Login authenticates a user and returns a JWT token
// @Summary User login
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginInput true "Login credentials"
// @Success 200 {object} TokenResponse "Successfully authenticated"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	ctx := context.Background()
	token, err := h.authSvc.Login(ctx, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken: token,
		TokenType:   "bearer",
		ExpiresIn:   3600, // 1 hour
	})
}

// Register creates a new user account
// @Summary Register a new user
// @Description Create a new user account with the provided information
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterInput true "User registration data"
// @Success 201 {object} user.User "Successfully registered user"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 409 {object} ErrorResponse "Email already exists"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	ctx := context.Background()
	user, err := h.authSvc.Register(ctx, input.Name, input.Email, input.Password)
	if err != nil {
		status := http.StatusInternalServerError
		errMsg := err.Error()

		// Handle specific error cases
		if errMsg == "email already exists" {
			status = http.StatusConflict
			errMsg = "Email is already registered"
		}

		c.JSON(status, ErrorResponse{Error: errMsg})
		return
	}

	// Clear sensitive data before sending response
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

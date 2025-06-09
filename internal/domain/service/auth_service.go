package service

import (
	"clean-arch-go/internal/domain/repository"
	"clean-arch-go/internal/domain/user"
	"clean-arch-go/internal/errors"
	"clean-arch-go/internal/pkg/redis"
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, name, email, password string) (*user.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	ValidateToken(tokenString string) (string, error)
	Logout(ctx context.Context, token string) error
	GetUserByID(ctx context.Context, id string) (*user.User, error)
	GenerateToken(ctx context.Context, userID string) (string, error)
	VerifyToken(ctx context.Context, token string) (string, error)
	RefreshToken(ctx context.Context, token string) (string, error)
	RevokeToken(ctx context.Context, token string) error
	GetUserFromToken(ctx context.Context, token string) (*user.User, error)
}

type authService struct {
	userRepo    repository.UserRepository
	tokenSecret string
	redisClient *redis.RedisClient
}

func NewAuthService(userRepo repository.UserRepository, tokenSecret string, redisClient *redis.RedisClient) AuthService {
	return &authService{
		userRepo:    userRepo,
		tokenSecret: tokenSecret,
		redisClient: redisClient,
	}
}

// Register creates a new user with the provided information
func (s *authService) Register(ctx context.Context, name, email, password string) (*user.User, error) {
	// Check if email already exists
	existingUser, _ := s.userRepo.FindByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.NewAppError("EMAIL_EXISTS", "Email already exists", nil)
	}

	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.NewAppError("PASSWORD_HASH_ERROR", "Failed to hash password", err)
	}

	// Create new user
	now := time.Now()
	user := &user.User{
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save user to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.NewAppError("USER_CREATION_ERROR", "Failed to create user", err)
	}

	// Clear password before returning
	user.Password = ""
	return user, nil
}

func (s *authService) GenerateToken(ctx context.Context, userID string) (string, error) {
	// TODO: Implement token generation
	return "", errors.NewAppError("NOT_IMPLEMENTED", "Token generation not implemented", nil)
}

func (s *authService) VerifyToken(ctx context.Context, token string) (string, error) {
	// TODO: Implement token verification
	return "", errors.NewAppError("NOT_IMPLEMENTED", "Token verification not implemented", nil)
}

func (s *authService) RefreshToken(ctx context.Context, token string) (string, error) {
	// TODO: Implement token refresh
	return "", errors.NewAppError("NOT_IMPLEMENTED", "Token refresh not implemented", nil)
}

func (s *authService) RevokeToken(ctx context.Context, token string) error {
	// TODO: Implement token revocation
	return errors.NewAppError("NOT_IMPLEMENTED", "Token revocation not implemented", nil)
}

// GetUserFromToken validates the JWT token and returns the associated user
func (s *authService) GetUserFromToken(ctx context.Context, tokenString string) (*user.User, error) {
	// First, check if the token is valid and get the user ID from the claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.NewAppError("INVALID_TOKEN", "Unexpected signing method", nil)
		}
		return []byte(s.tokenSecret), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.NewAppError("INVALID_TOKEN", "Malformed token", nil)
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, errors.NewAppError("TOKEN_EXPIRED", "Token is either expired or not active yet", nil)
			}
		}
		return nil, errors.NewAppError("INVALID_TOKEN", "Invalid token", err)
	}

	if !token.Valid {
		return nil, errors.NewAppError("INVALID_TOKEN", "Invalid token", nil)
	}

	// Extract the user ID from the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.NewAppError("INVALID_TOKEN", "Invalid token claims", nil)
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return nil, errors.NewAppError("INVALID_TOKEN", "Invalid user ID in token", nil)
	}

	// Check if the session exists in Redis
	sessionKey := "session:" + tokenString
	sessionUserID, err := s.redisClient.Get(ctx, sessionKey)
	if err != nil || sessionUserID != userID {
		return nil, errors.NewAppError("INVALID_SESSION", "Invalid or expired session", nil)
	}

	// Get the user from the database
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.NewAppError("USER_NOT_FOUND", "User not found", err)
	}

	if user == nil {
		return nil, errors.NewAppError("USER_NOT_FOUND", "User not found", nil)
	}

	// Clear the password before returning
	user.Password = ""
	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.NewAppError("NOT_FOUND", "User not found", nil)
	}

	// Check if user exists and is active
	if user == nil || user.ID == "" {
		return "", errors.NewAppError("INVALID_CREDENTIALS", "Invalid email or password", nil)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.NewAppError("INVALID_CREDENTIALS", "Invalid email or password", nil)
	}

	// Create JWT token with user claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":     time.Now().Unix(),                     // Issued at
	})


	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(s.tokenSecret))
	if err != nil {
		return "", errors.NewAppError("TOKEN_GENERATION_ERROR", "Failed to generate token", err)
	}

	// Store session in Redis with TTL
	sessionKey := "session:" + tokenString
	sessionTTL := 24 * time.Hour // 24 hours

	if err := s.redisClient.Set(ctx, sessionKey, user.ID, sessionTTL); err != nil {
		return "", errors.NewAppError("SESSION_STORAGE_ERROR", "Failed to store session", err)
	}

	return tokenString, nil
}

func (s *authService) ValidateToken(tokenString string) (string, error) {
	// Kiểm tra session trong Redis
	sessionKey := "session:" + tokenString
	userID, err := s.redisClient.Get(context.Background(), sessionKey)
	if err != nil {
		return "", errors.NewAppError("INVALID_SESSION", "Invalid or expired session", nil)
	}

	return userID, nil
}

// GetUserByID lấy thông tin user theo ID
func (s *authService) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// Ẩn mật khẩu trước khi trả về
	user.Password = ""
	return user, nil
}

// Logout invalidates a user session
func (s *authService) Logout(ctx context.Context, token string) error {
	sessionKey := "session:" + token
	return s.redisClient.Del(ctx, sessionKey)
}

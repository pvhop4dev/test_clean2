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

func (s *authService) Register(ctx context.Context, name, email, password string) (*user.User, error) {
	// Kiểm tra email đã tồn tại chưa
	existingUser, _ := s.userRepo.FindByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.NewAppError("EMAIL_EXISTS", "Email already exists", nil)
	}

	// Tạo user mới
	user := &user.User{
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

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

func (s *authService) GetUserFromToken(ctx context.Context, token string) (*user.User, error) {
	// TODO: Implement get user from token
	return nil, errors.NewAppError("NOT_IMPLEMENTED", "Get user from token not implemented", nil)
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	// Tìm user theo email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.NewAppError("NOT_FOUND", "User not found", nil)
	}

	// Kiểm tra mật khẩu
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.NewAppError("INVALID_CREDENTIALS", "Invalid credentials", nil)
	}

	// Tạo token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token hết hạn sau 24h
	})

	tokenString, err := token.SignedString([]byte(s.tokenSecret))
	if err != nil {
		return "", err
	}

	// Lưu session vào Redis
	sessionKey := "session:" + tokenString
	if err := s.redisClient.Set(ctx, sessionKey, user.ID, time.Hour*24); err != nil {
		return "", err
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

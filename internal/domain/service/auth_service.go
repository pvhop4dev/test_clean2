package service

import (
	"clean-arch-go/internal/domain/entities"
	"clean-arch-go/internal/domain/repository"
	"clean-arch-go/internal/pkg/config"
	"context"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, name, email, password string) (*entities.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	ValidateToken(tokenString string) (uint, error)
	GetUserByID(ctx context.Context, id uint) (*entities.User, error)
}

type authService struct {
	userRepo    repository.UserRepository
	tokenSecret string
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo:     userRepo,
		tokenSecret: cfg.App.Secret,
	}
}

func (s *authService) Register(ctx context.Context, name, email, password string) (*entities.User, error) {
	// Kiểm tra email đã tồn tại chưa
	existingUser, _ := s.userRepo.FindByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Mã hóa mật khẩu
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Tạo user mới
	user := &entities.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Ẩn mật khẩu trước khi trả về
	user.Password = ""
	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	// Tìm user theo email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Kiểm tra mật khẩu
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
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

	return tokenString, nil
}

func (s *authService) ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra phương thức mã hóa
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.tokenSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, errors.New("invalid token")
}

// GetUserByID lấy thông tin user theo ID
func (s *authService) GetUserByID(ctx context.Context, id uint) (*entities.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// Ẩn mật khẩu trước khi trả về
	user.Password = ""
	return user, nil
}

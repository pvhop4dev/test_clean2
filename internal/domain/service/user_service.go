package service

import (
	"clean-arch-go/internal/domain/repository"
	"clean-arch-go/internal/domain/user"
	"clean-arch-go/internal/errors"
	"context"
	"time"
)

type (
	User         = user.User
	TokenResponse = user.TokenResponse
	LoginInput    = user.LoginInput
	CreateUserInput = user.CreateUserInput
	UpdateUserInput = user.UpdateUserInput
)

type UserService interface {
	CreateUser(ctx context.Context, input *CreateUserInput) (*User, error)
	Login(ctx context.Context, input *LoginInput) (*TokenResponse, error)
	GetUser(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, input *UpdateUserInput) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	Register(ctx context.Context, name, email, password string) (*User, error)
}

type userService struct {
	userRepository repository.UserRepository
	authService    AuthService
}

func NewUserService(userRepository repository.UserRepository, authService AuthService) UserService {
	return &userService{
		userRepository: userRepository,
		authService:    authService,
	}
}

func (s *userService) CreateUser(ctx context.Context, input *user.CreateUserInput) (*user.User, error) {
	user := &user.User{
		Email:     input.Email,
		Password:  input.Password, // TODO: Hash password before saving
		Name:      input.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, input *LoginInput) (*TokenResponse, error) {
	user, err := s.userRepository.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.NewAppError("INVALID_CREDENTIALS", "Invalid email or password", nil)
	}

	// TODO: Verify password

	// TODO: Generate JWT token
	tokenResponse := &TokenResponse{
		Token: "TODO: Generate JWT token",
	}

	return tokenResponse, nil
}

func (s *userService) Register(ctx context.Context, name, email, password string) (*user.User, error) {
	return s.authService.Register(ctx, name, email, password)
}

func (s *userService) GetUser(ctx context.Context, id string) (*user.User, error) {
	return s.userRepository.FindByID(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, input *user.UpdateUserInput) (*user.User, error) {
	user, err := s.userRepository.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.NewAppError("NOT_FOUND", "User not found", nil)
	}

	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Name != "" {
		user.Name = input.Name
	}

	user.UpdatedAt = time.Now()

	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := s.userRepository.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepository.Delete(ctx, id)
}

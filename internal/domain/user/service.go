package user

import (
	"clean-arch-go/internal/errors"
	"context"
)

// Service defines the interface for user business logic
//go:generate mockery --name=Service --output=../mocks/user

type Service interface {
	CreateUser(ctx context.Context, input CreateUserInput) (*User, error)
	GetUser(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, id string, input UpdateUserInput) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, page, limit int) ([]*User, error)
}

// NewService creates a new user service
func NewService(repo Repository) Service {
	return &userService{
		repo: repo,
	}
}

type userService struct {
	repo Repository
}

func (s *userService) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	existingUser, err := s.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.NewAppError("DUPLICATE_EMAIL", "Email already exists", nil)
	}

	user := &User{
		Email:    input.Email,
		Password: input.Password, // In real app, hash password here
		Name:     input.Name,
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, user)
}

func (s *userService) GetUser(ctx context.Context, id string) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, id string, input UpdateUserInput) (*User, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	existingUser, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:    existingUser.ID,
		Email: input.Email,
		Name:  input.Name,
	}

	return s.repo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context, page, limit int) ([]*User, error) {
	return s.repo.List(ctx, page, limit)
}

package user

import (
	"clean-arch-go/internal/domain/repository"
	"clean-arch-go/internal/domain/user"
	"clean-arch-go/internal/errors"
	"context"
)

//go:generate mockery --name=UserUsecase --output=../mocks/user

type UserUsecase interface {
	CreateUser(ctx context.Context, input user.CreateUserInput) (*user.User, error)
	GetUser(ctx context.Context, id string) (*user.User, error)
	UpdateUser(ctx context.Context, id string, input user.UpdateUserInput) (*user.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, page, limit int) ([]*user.User, error)
	CountUsers(ctx context.Context) (int64, error)
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{
		repo: repo,
	}
}

func (u *userUsecase) CreateUser(ctx context.Context, input user.CreateUserInput) (*user.User, error) {
	if err := input.Validate(); err != nil {
		return nil, errors.NewValidationError("input", err.Error())
	}

	user := &user.User{
		Email:    input.Email,
		Password: input.Password,
		Name:     input.Name,
	}

	if err := user.Validate(); err != nil {
		return nil, errors.NewValidationError("user", err.Error())
	}

	if err := u.repo.Create(ctx, user); err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}

	return user, nil
}

func (u *userUsecase) GetUser(ctx context.Context, id string) (*user.User, error) {
	result, err := u.repo.FindByID(ctx, id)
	if err != nil {
		if err == user.ErrUserNotFound {
			return nil, errors.NewNotFoundError("user")
		}
		return nil, errors.NewInternalServerError(err.Error())
	}

	return result, nil
}

func (u *userUsecase) UpdateUser(ctx context.Context, id string, input user.UpdateUserInput) (*user.User, error) {
	if err := input.Validate(); err != nil {
		return nil, errors.NewValidationError("input", err.Error())
	}

	existingUser, err := u.repo.FindByID(ctx, id)
	if err != nil {
		if err == user.ErrUserNotFound {
			return nil, errors.NewNotFoundError("user")
		}
		return nil, errors.NewInternalServerError(err.Error())
	}

	updatedUser := &user.User{
		ID:        id,
		Email:     input.Email,
		Name:      input.Name,
		UpdatedAt: existingUser.UpdatedAt,
	}

	if err := u.repo.Update(ctx, updatedUser); err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}

	return updatedUser, nil
}

func (u *userUsecase) DeleteUser(ctx context.Context, id string) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		if err == user.ErrUserNotFound {
			return errors.NewNotFoundError("user")
		}
		return errors.NewInternalServerError(err.Error())
	}

	return nil
}

func (u *userUsecase) ListUsers(ctx context.Context, page, limit int) ([]*user.User, error) {
	users, err := u.repo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	return users, nil
}

func (u *userUsecase) CountUsers(ctx context.Context) (int64, error) {
	count, err := u.repo.Count(ctx, nil)
	if err != nil {
		return 0, errors.NewInternalServerError(err.Error())
	}
	return count, nil
}

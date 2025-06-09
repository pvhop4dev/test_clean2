package user

import (
	"clean-arch-go/internal/errors"
	"context"

	"gorm.io/gorm"
)

// Repository defines the interface for user data access
//go:generate mockery --name=Repository --output=../mocks/user

type Repository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}

// ErrUserNotFound is returned when user is not found
var ErrUserNotFound = errors.NewAppError("USER_NOT_FOUND", "User not found", nil)

// NewRepository creates a new user repository
func NewRepository(db *gorm.DB) Repository {
	return &userRepository{
		db: db,
	}
}

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) Create(ctx context.Context, user *User) (*User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*User, error) {
	var user User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, errors.NewInternalServerError(err.Error())
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *User) (*User, error) {
	if err := r.db.Save(user).Error; err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.Delete(&User{}, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrUserNotFound
		}
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}

func (r *userRepository) List(ctx context.Context, page, limit int) ([]*User, error) {
	var users []*User
	if err := r.db.Offset((page - 1) * limit).Limit(limit).Find(&users).Error; err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	return users, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.NewInternalServerError(err.Error())
	}
	return &user, nil
}

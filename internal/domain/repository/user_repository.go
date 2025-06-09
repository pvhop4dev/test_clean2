package repository

import (
	"clean-arch-go/internal/domain/user"
	"clean-arch-go/internal/errors"
	"context"

	"gorm.io/gorm"
	"clean-arch-go/internal/pkg/database"
)

type UserRepository interface {
	BaseRepository[user.User]
	FindByEmail(ctx context.Context, email string) (*user.User, error)
}

type userRepository struct {
	*baseRepository[user.User]
}

func NewUserRepository(db *database.Database) UserRepository {
	return &userRepository{
		baseRepository: NewBaseRepository[user.User](db.DB).(*baseRepository[user.User]),
	}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var user user.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.NewInternalServerError(err.Error())
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	var user user.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&user.User{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Count(ctx context.Context, query interface{}) (int64, error) {
	return r.baseRepository.Count(ctx, query)
}

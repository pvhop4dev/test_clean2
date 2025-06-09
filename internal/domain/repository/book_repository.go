package repository

import (
	"clean-arch-go/internal/domain/entities"
	"clean-arch-go/internal/errors"
	"context"

	"clean-arch-go/internal/pkg/database"
)

type BookRepository interface {
	BaseRepository[entities.Book]
	ListByUserID(ctx context.Context, userID string, page, limit int) ([]*entities.Book, error)
}

type bookRepository struct {
	*baseRepository[entities.Book]
}

func NewBookRepository(db *database.Database) BookRepository {
	return &bookRepository{
		baseRepository: NewBaseRepository[entities.Book](db.DB).(*baseRepository[entities.Book]),
	}
}

func (r *bookRepository) ListByUserID(ctx context.Context, userID string, page, limit int) ([]*entities.Book, error) {
	var books []*entities.Book
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&books).Error; err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	return books, nil
}

func (r *bookRepository) Count(ctx context.Context, query interface{}) (int64, error) {
	return r.baseRepository.Count(ctx, query)
}

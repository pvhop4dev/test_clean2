package repository

import (
	"clean-arch-go/internal/domain/entities"
	"context"

	"gorm.io/gorm"
)

type BookRepository interface {
	BaseRepository[entities.Book]
	FindByUserID(ctx context.Context, userID uint) ([]*entities.Book, error)
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(ctx context.Context, book *entities.Book) error {
	return r.db.WithContext(ctx).Create(book).Error
}

func (r *bookRepository) Update(ctx context.Context, book *entities.Book) error {
	return r.db.WithContext(ctx).Save(book).Error
}

func (r *bookRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.Book{}, id).Error
}

func (r *bookRepository) FindByID(ctx context.Context, id uint) (*entities.Book, error) {
	var book entities.Book
	err := r.db.WithContext(ctx).First(&book, id).Error
	return &book, err
}

func (r *bookRepository) FindByUserID(ctx context.Context, userID uint) ([]*entities.Book, error) {
	var books []*entities.Book
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&books).Error
	return books, err
}

func (r *bookRepository) FindAll(ctx context.Context) ([]*entities.Book, error) {
	var books []*entities.Book
	err := r.db.WithContext(ctx).Find(&books).Error
	return books, err
}

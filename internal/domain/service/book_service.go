package service

import (
	"clean-arch-go/internal/domain/entities"
	"clean-arch-go/internal/domain/repository"
	"clean-arch-go/internal/errors"
	"context"
)

type BookService interface {
	CreateBook(ctx context.Context, book *entities.Book) error
	UpdateBook(ctx context.Context, id string, book *entities.Book) error
	DeleteBook(ctx context.Context, id string) error
	GetBookByID(ctx context.Context, id string) (*entities.Book, error)
	ListBooksByUserID(ctx context.Context, userID string, page, limit int) ([]*entities.Book, error)
	CheckBookOwnership(ctx context.Context, bookID, userID string) error
}

type bookService struct {
	bookRepo repository.BookRepository
}

func NewBookService(bookRepo repository.BookRepository) BookService {
	return &bookService{
		bookRepo: bookRepo,
	}
}

func (s *bookService) CreateBook(ctx context.Context, book *entities.Book) error {
	return s.bookRepo.Create(ctx, book)
}

func (s *bookService) UpdateBook(ctx context.Context, id string, book *entities.Book) error {
	// Kiểm tra xem sách có tồn tại không
	existingBook, err := s.bookRepo.FindByID(ctx, id)
	if err != nil {
		return errors.NewAppError("NOT_FOUND", "Book not found", nil)
	}

	// Cập nhật thông tin sách
	existingBook.Title = book.Title
	existingBook.Description = book.Description
	existingBook.Author = book.Author

	return s.bookRepo.Update(ctx, existingBook)
}

func (s *bookService) DeleteBook(ctx context.Context, id string) error {
	// Kiểm tra xem sách có tồn tại không
	book, err := s.bookRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if book == nil {
		return errors.NewAppError("NOT_FOUND", "Book not found", nil)
	}

	return s.bookRepo.Delete(ctx, id)
}

func (s *bookService) GetBookByID(ctx context.Context, id string) (*entities.Book, error) {
	book, err := s.bookRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if book == nil {
		return nil, errors.NewAppError("NOT_FOUND", "Book not found", nil)
	}

	return book, nil
}

func (s *bookService) CheckBookOwnership(ctx context.Context, bookID, userID string) error {
	book, err := s.bookRepo.FindByID(ctx, bookID)
	if err != nil {
		return err
	}

	if book == nil {
		return errors.NewAppError("NOT_FOUND", "Book not found", nil)
	}

	if book.UserID != userID {
		return errors.NewAppError("UNAUTHORIZED", "You are not authorized to access this book", nil)
	}

	return nil
}

func (s *bookService) ListBooksByUserID(ctx context.Context, userID string, page, limit int) ([]*entities.Book, error) {
	return s.bookRepo.ListByUserID(ctx, userID, page, limit)
}

package service

import (
	"clean-arch-go/internal/domain/entities"
	"clean-arch-go/internal/domain/repository"
	"context"
	"errors"
)

type BookService interface {
	CreateBook(ctx context.Context, book *entities.Book) error
	UpdateBook(ctx context.Context, book *entities.Book) error
	DeleteBook(ctx context.Context, id, userID uint) error
	GetBookByID(ctx context.Context, id, userID uint) (*entities.Book, error)
	ListBooksByUserID(ctx context.Context, userID uint) ([]*entities.Book, error)
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

func (s *bookService) UpdateBook(ctx context.Context, book *entities.Book) error {
	// Kiểm tra xem sách có tồn tại không
	existingBook, err := s.bookRepo.FindByID(ctx, book.ID)
	if err != nil {
		return errors.New("book not found")
	}

	// Cập nhật thông tin sách
	existingBook.Title = book.Title
	existingBook.Description = book.Description
	existingBook.Author = book.Author

	return s.bookRepo.Update(ctx, existingBook)
}

func (s *bookService) DeleteBook(ctx context.Context, id, userID uint) error {
	// Kiểm tra xem sách có tồn tại và thuộc về user không
	book, err := s.bookRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("book not found")
	}

	if book.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.bookRepo.Delete(ctx, id)
}

func (s *bookService) GetBookByID(ctx context.Context, id, userID uint) (*entities.Book, error) {
	book, err := s.bookRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("book not found")
	}

	// Chỉ trả về sách nếu nó thuộc về user hoặc có cơ chế phân quyền phù hợp
	if book.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	return book, nil
}

func (s *bookService) ListBooksByUserID(ctx context.Context, userID uint) ([]*entities.Book, error) {
	return s.bookRepo.FindByUserID(ctx, userID)
}

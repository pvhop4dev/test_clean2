package book

import (
	"context"
	"errors"
)

var (
	// ErrBookNotFound is returned when a book is not found
	ErrBookNotFound = errors.New("book not found")
)

// Service handles business logic for books
type Service interface {
	CreateBook(ctx context.Context, book *Book) error
	GetBook(ctx context.Context, id string) (*Book, error)
	UpdateBook(ctx context.Context, book *Book) error
	DeleteBook(ctx context.Context, id string) error
	ListBooks(ctx context.Context, filter Filter) ([]*Book, error)
}

type service struct {
	repo Repository
}

// NewService creates a new book service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateBook(ctx context.Context, book *Book) error {
	return s.repo.Create(ctx, book)
}

func (s *service) GetBook(ctx context.Context, id string) (*Book, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) UpdateBook(ctx context.Context, book *Book) error {
	return s.repo.Update(ctx, book)
}

func (s *service) DeleteBook(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) ListBooks(ctx context.Context, filter Filter) ([]*Book, error) {
	return s.repo.FindAll(ctx, filter)
}

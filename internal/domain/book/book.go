package book

import "context"

// Book represents a book in the system
type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	// Add other book fields as needed
}

// Repository defines the interface for book data operations
type Repository interface {
	Create(ctx context.Context, book *Book) error
	Update(ctx context.Context, book *Book) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*Book, error)
	FindAll(ctx context.Context, filter Filter) ([]*Book, error)
}

// Filter represents the filter criteria for querying books
type Filter struct {
	Title  string `json:"title,omitempty"`
	Author string `json:"author,omitempty"`
	// Add other filter fields as needed
}

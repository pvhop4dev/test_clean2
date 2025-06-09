package book

import "context"

// Define a key type to avoid key collisions in context
type contextKey string

const (
	// BookCtxKey is the key used to store book in context
	BookCtxKey contextKey = "book"
)

// NewContextWithBook returns a new context with the book added to it
func NewContextWithBook(ctx context.Context, book *Book) context.Context {
	return context.WithValue(ctx, BookCtxKey, book)
}

// BookFromContext retrieves the book from the context if it exists
func BookFromContext(ctx context.Context) (*Book, bool) {
	book, ok := ctx.Value(BookCtxKey).(*Book)
	return book, ok
}

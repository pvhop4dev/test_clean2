package cached

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"clean-arch-go/internal/domain/entities"
	"clean-arch-go/internal/domain/repository"
	"clean-arch-go/internal/pkg/redis"
)

type cachedBookRepository struct {
	repo    repository.BookRepository
	cache   *redis.RedisClient
	keyFunc func(id string) string
}

// NewCachedBookRepository creates a new cached book repository
func NewCachedBookRepository(
	repo repository.BookRepository,
	cache *redis.RedisClient,
) repository.BookRepository {
	keyFunc := func(id string) string {
		return fmt.Sprintf("book:%s", id)
	}

	return &cachedBookRepository{
		repo:    repo,
		cache:   cache,
		keyFunc: keyFunc,
	}
}

// Helper methods for caching
func (r *cachedBookRepository) getCached(ctx context.Context, key string, dest *entities.Book) (bool, error) {
	data, err := r.cache.Get(ctx, key)
	if err != nil || data == "" {
		return false, err
	}
	return true, json.Unmarshal([]byte(data), dest)
}

func (r *cachedBookRepository) setCached(ctx context.Context, key string, value entities.Book) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal book: %w", err)
	}
	return r.cache.Set(ctx, key, string(data), 5*time.Minute)
}

func (r *cachedBookRepository) deleteCached(ctx context.Context, key string) error {
	return r.cache.Del(ctx, key)
}

func (r *cachedBookRepository) FindByID(ctx context.Context, id string) (*entities.Book, error) {
	key := r.keyFunc(id)

	// Try to get from cache
	var b entities.Book
	found, err := r.getCached(ctx, key, &b)
	if err != nil {
		return nil, err
	}
	if found {
		return &b, nil
	}

	// Not in cache, get from repository
	book, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if book != nil {
		if err := r.setCached(ctx, key, *book); err != nil {
			// Log error but don't fail the request
			log.Printf("Failed to cache book %s: %v", id, err)
		}
		return book, nil
	}

	return nil, nil
}

func (r *cachedBookRepository) Create(ctx context.Context, b *entities.Book) error {
	if err := r.repo.Create(ctx, b); err != nil {
		return err
	}

	// Invalidate cache
	key := r.keyFunc(b.ID)
	return r.deleteCached(ctx, key)
}

func (r *cachedBookRepository) Update(ctx context.Context, b *entities.Book) error {
	if err := r.repo.Update(ctx, b); err != nil {
		return err
	}

	// Invalidate cache
	key := r.keyFunc(b.ID)
	return r.deleteCached(ctx, key)
}

// Delete deletes a book by ID and invalidates the cache
func (r *cachedBookRepository) Delete(ctx context.Context, id string) error {
	// Invalidate cache first to avoid race conditions
	key := r.keyFunc(id)
	if err := r.deleteCached(ctx, key); err != nil {
		log.Printf("Failed to invalidate cache for book %s: %v", id, err)
	}

	// Delete from repository
	return r.repo.Delete(ctx, id)
}

// Count returns the count of books matching the query
func (r *cachedBookRepository) Count(ctx context.Context, query interface{}) (int64, error) {
	return r.repo.Count(ctx, query)
}

// FindAll returns all books with pagination
func (r *cachedBookRepository) FindAll(ctx context.Context, page, limit int) ([]*entities.Book, error) {
	return r.repo.FindAll(ctx, page, limit)
}

// FindOne finds a single book based on the query
func (r *cachedBookRepository) FindOne(ctx context.Context, query interface{}) (*entities.Book, error) {
	return r.repo.FindOne(ctx, query)
}

// FindMany finds multiple books based on the query with pagination
func (r *cachedBookRepository) FindMany(ctx context.Context, query interface{}, page, limit int) ([]*entities.Book, error) {
	return r.repo.FindMany(ctx, query, page, limit)
}

// ListByUserID lists books by user ID with pagination
func (r *cachedBookRepository) ListByUserID(ctx context.Context, userID string, page, limit int) ([]*entities.Book, error) {
	return r.repo.ListByUserID(ctx, userID, page, limit)
}

package repository

import (
	"context"
	"fmt"

	"clean-arch-go/internal/domain/book"
	"clean-arch-go/internal/domain/entities"
	"clean-arch-go/internal/pkg/cache"
)

const (
	bookCacheTTL = 300 // 5 minutes in seconds
	bookCacheKey = "book:%s"
)

type cachedBookRepository struct {
	repo  book.Repository
	cache cache.Cache
}

func NewCachedBookRepository(repo book.Repository, cache cache.Cache) book.Repository {
	return &cachedBookRepository{
		repo:  repo,
		cache: cache,
	}
}

func (r *cachedBookRepository) FindByID(ctx context.Context, id string) (*entities.Book, error) {
	cacheKey := fmt.Sprintf(bookCacheKey, id)

	var cachedBook entities.Book
	found, err := r.cache.Get(ctx, cacheKey, &cachedBook)
	if err != nil {
		return nil, err
	}
	if found {
		return &cachedBook, nil
	}

	b, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if b != nil {
		_ = r.cache.Set(ctx, cacheKey, b, bookCacheTTL)
	}

	return b, nil
}

func (r *cachedBookRepository) FindAll(ctx context.Context, filter book.Filter) ([]*entities.Book, error) {
	// For simplicity, we'll go to the database directly for filtered queries
	return r.repo.FindAll(ctx, filter)
}

func (r *cachedBookRepository) Create(ctx context.Context, book *entities.Book) error {
	if err := r.repo.Create(ctx, book); err != nil {
		return err
	}
	return nil
}

func (r *cachedBookRepository) Update(ctx context.Context, book *entities.Book) error {
	if err := r.repo.Update(ctx, book); err != nil {
		return err
	}

	cacheKey := fmt.Sprintf(bookCacheKey, book.ID)
	return r.cache.Delete(ctx, cacheKey)
}

func (r *cachedBookRepository) Delete(ctx context.Context, id string) error {
	if err := r.repo.Delete(ctx, id); err != nil {
		return err
	}

	cacheKey := fmt.Sprintf(bookCacheKey, id)
	return r.cache.Delete(ctx, cacheKey)
}

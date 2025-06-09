// internal/infrastructure/repository/cached_user_repository.go
package repository

import (
	"context"
	"fmt"
	"log"

	"clean-arch-go/internal/domain/user"
	"clean-arch-go/internal/pkg/cache"
)

const (
	userCacheTTL = 300 // 5 minutes in seconds
	userCacheKey = "user:%s"
)

type cachedUserRepository struct {
	repo  user.Repository
	cache cache.Cache
}

func NewCachedUserRepository(repo user.Repository, cache cache.Cache) user.Repository {
	return &cachedUserRepository{
		repo:  repo,
		cache: cache,
	}
}

func (r *cachedUserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	cacheKey := fmt.Sprintf(userCacheKey, id)

	var cachedUser user.User
	found, err := r.cache.Get(ctx, cacheKey, &cachedUser)
	if err != nil {
		return nil, err
	}
	if found {
		return &cachedUser, nil
	}

	// If not in cache, get from repository
	u, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if err := r.cache.Set(ctx, cacheKey, u, userCacheTTL); err != nil {
		// Log error but don't fail the request
		log.Printf("Failed to cache user %s: %v", id, err)
	}

	return u, nil
}

func (r *cachedUserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	// For GetByID, we can just use FindByID since it has the same signature
	return r.FindByID(ctx, id)
}

func (r *cachedUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	// Note: We're not caching by email as it's better to cache by ID
	return r.repo.FindByEmail(ctx, email)
}

func (r *cachedUserRepository) List(ctx context.Context, page, limit int) ([]*user.User, error) {
	// For listing users, we typically don't cache the results as they can change frequently
	// and the cache invalidation would be complex
	return r.repo.List(ctx, page, limit)
}

func (r *cachedUserRepository) Update(ctx context.Context, user *user.User) (*user.User, error) {
	// Invalidate cache
	cacheKey := fmt.Sprintf(userCacheKey, user.ID)
	if err := r.cache.Delete(ctx, cacheKey); err != nil {
		return nil, err
	}

	// Update in repository
	return r.repo.Update(ctx, user)
}

func (r *cachedUserRepository) Delete(ctx context.Context, id string) error {
	if err := r.repo.Delete(ctx, id); err != nil {
		return err
	}

	cacheKey := fmt.Sprintf(userCacheKey, id)
	return r.cache.Delete(ctx, cacheKey)
}

func (r *cachedUserRepository) Create(ctx context.Context, user *user.User) (*user.User, error) {
	u, err := r.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf(userCacheKey, user.ID)
	_ = r.cache.Delete(ctx, cacheKey)

	return u, nil
}

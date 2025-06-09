package cached

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"clean-arch-go/internal/domain/repository"
	"clean-arch-go/internal/domain/user"
	"clean-arch-go/internal/pkg/redis"
)

// cachedUserRepository implements a cached version of UserRepository
type cachedUserRepository struct {
	repo    repository.UserRepository
	cache   *redis.RedisClient
	keyFunc func(id string) string
}

// NewCachedUserRepository creates a new cached user repository
func NewCachedUserRepository(
	repo repository.UserRepository,
	cache *redis.RedisClient,
) repository.UserRepository {
	keyFunc := func(id string) string {
		return fmt.Sprintf("user:%s", id)
	}

	return &cachedUserRepository{
		repo:    repo,
		cache:   cache,
		keyFunc: keyFunc,
	}
}

// FindByEmail finds a user by email (not cached)
func (r *cachedUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	// For email-based lookups, we don't cache by default
	return r.repo.FindByEmail(ctx, email)
}

// FindByID finds a user by ID with caching
func (r *cachedUserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	key := r.keyFunc(id)

	// Try to get from cache
	var u user.User
	found, err := r.getCached(ctx, key, &u)
	if err != nil {
		log.Printf("Error getting from cache: %v", err)
		// Continue with database lookup on cache error
	}

	if found {
		return &u, nil
	}

	// Not in cache, get from repository
	user, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if user != nil {
		if err := r.setCached(ctx, key, *user); err != nil {
			log.Printf("Failed to cache user %s: %v", id, err)
		}
	}

	return user, nil
}

// Create creates a new user and invalidates cache
func (r *cachedUserRepository) Create(ctx context.Context, u *user.User) error {
	if err := r.repo.Create(ctx, u); err != nil {
		return err
	}

	// Invalidate cache
	key := r.keyFunc(u.ID)
	return r.deleteCached(ctx, key)
}

// Update updates a user and invalidates cache
func (r *cachedUserRepository) Update(ctx context.Context, u *user.User) error {
	if err := r.repo.Update(ctx, u); err != nil {
		return err
	}

	// Invalidate cache
	key := r.keyFunc(u.ID)
	return r.deleteCached(ctx, key)
}

// Delete deletes a user and invalidates cache
func (r *cachedUserRepository) Delete(ctx context.Context, id string) error {
	// Get user first to invalidate cache
	user, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete from repository
	if err := r.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Invalidate cache
	key := r.keyFunc(user.ID)
	return r.deleteCached(ctx, key)
}

// Count returns the count of users matching the query
func (r *cachedUserRepository) Count(ctx context.Context, query interface{}) (int64, error) {
	// For count operations, we don't use cache
	return r.repo.Count(ctx, query)
}

// FindAll returns all users with pagination
func (r *cachedUserRepository) FindAll(ctx context.Context, page, limit int) ([]*user.User, error) {
	// For listing operations, we don't use cache by default
	return r.repo.FindAll(ctx, page, limit)
}

// FindMany finds multiple users based on the query with pagination
func (r *cachedUserRepository) FindMany(ctx context.Context, query interface{}, page, limit int) ([]*user.User, error) {
	// For query operations, we don't use cache by default
	return r.repo.FindMany(ctx, query, page, limit)
}

// FindOne finds a single user based on the query
func (r *cachedUserRepository) FindOne(ctx context.Context, query interface{}) (*user.User, error) {
	// For single query operations, we don't use cache by default
	return r.repo.FindOne(ctx, query)
}

// Helper methods that delegate to the base repository
func (r *cachedUserRepository) getCached(ctx context.Context, key string, dest *user.User) (bool, error) {
	if r.cache == nil {
		return false, nil
	}

	data, err := r.cache.Get(ctx, key)
	if err != nil {
		return false, err
	}
	if data == "" {
		return false, nil
	}

	return true, json.Unmarshal([]byte(data), dest)
}

func (r *cachedUserRepository) setCached(ctx context.Context, key string, value user.User) error {
	if r.cache == nil {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.cache.Set(ctx, key, string(data), 5*time.Minute)
}

func (r *cachedUserRepository) deleteCached(ctx context.Context, key string) error {
	if r.cache == nil {
		return nil
	}
	return r.cache.Del(ctx, key)
}

package cached

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"clean-arch-go/internal/domain/repository"
	"clean-arch-go/internal/pkg/redis"
)

// baseCachedRepository provides common caching functionality for repositories
type baseCachedRepository[E any, R repository.BaseRepository[E]] struct {
	repo    R
	cache   *redis.RedisClient
	keyFunc func(id string) string
	ttl     time.Duration
}

// newBaseCachedRepository creates a new base cached repository
func newBaseCachedRepository[E any, R repository.BaseRepository[E]](
	repo R,
	cache *redis.RedisClient,
	keyFunc func(id string) string,
) *baseCachedRepository[E, R] {
	return &baseCachedRepository[E, R]{
		repo:    repo,
		cache:   cache,
		keyFunc: keyFunc,
		ttl:     5 * time.Minute, // Default TTL
	}
}

// getCached retrieves a value from cache
func (r *baseCachedRepository[E, R]) getCached(ctx context.Context, key string, dest *E) (bool, error) {
	if r.cache == nil {
		return false, nil
	}

	data, err := r.cache.Get(ctx, key)
	if err != nil {
		log.Printf("Error getting from cache: %v", err)
		return false, nil
	}
	if data == "" {
		return false, nil
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		log.Printf("Error unmarshaling cached data: %v", err)
		return false, nil
	}

	return true, nil
}

// setCached stores a value in the cache
func (r *baseCachedRepository[E, R]) setCached(ctx context.Context, key string, value E) error {
	if r.cache == nil {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.cache.Set(ctx, key, string(data), r.ttl)
}

// deleteCached removes a value from the cache
func (r *baseCachedRepository[E, R]) deleteCached(ctx context.Context, key string) error {
	if r.cache == nil {
		return nil
	}

	return r.cache.Del(ctx, key)
}

// BaseRepository is the base interface that all cached repositories must implement
type BaseRepository[E any] interface {
	repository.BaseRepository[E]
	getCached(ctx context.Context, key string, dest *E) (bool, error)
	setCached(ctx context.Context, key string, value E) error
	deleteCached(ctx context.Context, key string) error
}

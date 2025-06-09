package cached

import (
	"context"
	"encoding/json"
	"time"

	"clean-arch-go/internal/domain/repository"
	"clean-arch-go/internal/pkg/redis"
)

type cachedRepository[E any, R repository.BaseRepository[E]] struct {
	repo    R
	cache   *redis.RedisClient
	keyFunc func(id string) string
	ttl     time.Duration
}

func newCachedRepository[E any, R repository.BaseRepository[E]](
	repo R,
	cache *redis.RedisClient,
	keyFunc func(id string) string,
) *cachedRepository[E, R] {
	return &cachedRepository[E, R]{
		repo:    repo,
		cache:   cache,
		keyFunc: keyFunc,
		ttl:     5 * time.Minute, // Default TTL
	}
}

func (r *cachedRepository[E, R]) getCached(ctx context.Context, key string, dest *E) (bool, error) {
	data, err := r.cache.Get(ctx, key)
	if err != nil {
		return false, err
	}
	if data == "" {
		return false, nil
	}
	return true, json.Unmarshal([]byte(data), dest)
}

func (r *cachedRepository[E, R]) setCached(ctx context.Context, key string, value E) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.cache.Set(ctx, key, string(data), r.ttl)
}

func (r *cachedRepository[E, R]) deleteCached(ctx context.Context, key string) error {
	return r.cache.Del(ctx, key)
}

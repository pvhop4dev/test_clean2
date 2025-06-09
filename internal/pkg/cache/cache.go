package cache

import "context"

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration int64) error
	Get(ctx context.Context, key string, dest interface{}) (bool, error)
	Delete(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
}

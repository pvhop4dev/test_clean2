// internal/pkg/cache/redis_cache.go
package cache

import (
    "context"
    "encoding/json"
    "time"

    "github.com/redis/go-redis/v9"
)

type redisCache struct {
    client *redis.Client
}

func NewRedisCache(client *redis.Client) *redisCache {
    return &redisCache{client: client}
}

func (r *redisCache) Set(ctx context.Context, key string, value interface{}, expiration int64) error {
    jsonData, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return r.client.Set(ctx, key, jsonData, time.Duration(expiration)*time.Second).Err()
}

func (r *redisCache) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
    data, err := r.client.Get(ctx, key).Bytes()
    if err == redis.Nil {
        return false, nil
    }
    if err != nil {
        return false, err
    }
    return true, json.Unmarshal(data, dest)
}

func (r *redisCache) Delete(ctx context.Context, keys ...string) error {
    return r.client.Del(ctx, keys...).Err()
}

func (r *redisCache) Exists(ctx context.Context, key string) (bool, error) {
    n, err := r.client.Exists(ctx, key).Result()
    return n > 0, err
}

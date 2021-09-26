package cache

import (
	"github.com/go-redis/redis/v8"
	"time"
)

type Cache struct {
	Client *redis.Client
}

func (cache *Cache) Get(key string) (string, error) {
	return cache.Client.Get(ctx, key).Result()
}

func (cache *Cache) Set(key, value string, expiration time.Duration) error {
	return cache.Client.Set(ctx, key, value, expiration).Err()
}

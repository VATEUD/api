package cache

import (
	"api/utils"
	"context"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var RedisCache *Cache

func New() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     utils.Getenv("REDIS_ADDRESS", "localhost:6379"),
		Password: utils.Getenv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	RedisCache = &Cache{rdb}
}


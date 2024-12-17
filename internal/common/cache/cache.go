package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewClient(cfg RedisConfig) (*redis.Client, error) {
	addr := cfg.Host + ":" + cfg.Port
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Clean up the context to avoid resource leaks

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}

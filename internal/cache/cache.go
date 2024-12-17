package cache

import (
	"context"
	"time"

	"github.com/ozzy-cox/automatic-message-system/config"
	"github.com/redis/go-redis/v9"
)

var CacheClient *redis.Client

func GetClient(cfg config.RedisConfig) (*redis.Client, error) {
	if CacheClient != nil {
		return CacheClient, nil
	}
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
	CacheClient = rdb

	return rdb, nil
}

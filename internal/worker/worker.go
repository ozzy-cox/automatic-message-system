package worker

import (
	"context"
	"strconv"
	"time"

	"github.com/ozzy-cox/automatic-message-system/internal/cache"
	"github.com/redis/go-redis/v9"
)

const offsetKey = "producer_offset"

func mustGetProducerOffset() int {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	offset, err := cache.CacheClient.Get(ctx, offsetKey).Result()
	if err != nil {
		if err == redis.Nil {
			return 0
		}
		panic("Can't get producer offset from redis")
	}
	offsetValue, err := strconv.Atoi(offset)

	if err != nil {
		panic("Can't get producer offset from redis")
	}

	return offsetValue
}

func mustSetProducerOffset(offsetValue *int) {
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Second)
	defer cancel()

	_, err := cache.CacheClient.Set(ctx, offsetKey, *offsetValue, redis.KeepTTL).Result()
	if err != nil {
		panic("Can't set producer offset to redis")
	}
}

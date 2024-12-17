package worker

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	dbConn      *sql.DB
	cacheClient *redis.Client
)

func Initialize(_dbConn *sql.DB, _cacheClient *redis.Client) {
	dbConn = _dbConn
	cacheClient = _cacheClient
}

const offsetKey = "producer_offset"

func mustGetProducerOffset() int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	offset, err := cacheClient.Get(ctx, offsetKey).Result()
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

	_, err := cacheClient.Set(ctx, offsetKey, *offsetValue, redis.KeepTTL).Result()
	if err != nil {
		panic("Can't set producer offset to redis")
	}
}

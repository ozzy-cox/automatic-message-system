package producer

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const offsetKey = "producer_offset"

func (s *Service) mustGetProducerOffset() int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	offset, err := s.Cache.Get(ctx, offsetKey).Result()
	if err != nil {
		if err == redis.Nil {
			s.Logger.Println("No offset found in cache, starting from 0")
			return 0
		}
		s.Logger.Fatalf("Failed to get producer offset from redis: %v", err)
	}
	offsetValue, err := strconv.Atoi(offset)

	if err != nil {
		s.Logger.Fatalf("Failed to parse offset value from redis: %v", err)
	}

	return offsetValue
}

func (s *Service) mustSetProducerOffset(offsetValue *int) {
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Second)
	defer cancel()

	_, err := s.Cache.Set(ctx, offsetKey, *offsetValue, redis.KeepTTL).Result()
	if err != nil {
		s.Logger.Fatalf("Failed to set producer offset in redis: %v", err)
	}
}

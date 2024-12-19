package consumer

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CachedMessageResponse struct {
	MessageResponse
	timestamp time.Time
}

func (s *Service) cacheMessageResponse(ctx context.Context, msgId string, msg MessageResponse) error {
	value := CachedMessageResponse{
		MessageResponse: msg,
		timestamp:       time.Now(),
	}
	_, err := s.Cache.Set(ctx, msgId, value, redis.KeepTTL).Result()
	if err != nil {
		s.Logger.Fatalf("Failed to cache message in redis: %v", err)
	}
	return nil
}

package consumer

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func (s *Service) cacheMessageResponse(ctx context.Context, msgId string, msg MessageResponse) error {
	_, err := s.Cache.Set(ctx, msgId, msg, redis.KeepTTL).Result()
	if err != nil {
		s.Logger.Fatalf("Failed to cache message in redis: %v", err)
	}
	return nil
}

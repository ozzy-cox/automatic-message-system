package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *Service) cacheMessageResponse(ctx context.Context, msgId string, msg MessageResponse) error {
	value := map[string]any{
		"message":   msg.Message,
		"messageId": msg.MessageId,
		"timestamp": time.Now(),
	}
	jsonData, err := json.Marshal(value)
	_, err = s.Cache.Set(ctx, msgId, jsonData, redis.KeepTTL).Result()
	if err != nil {
		s.Logger.Fatalf("Failed to cache message in redis: %v", err)
	}
	return nil
}

package producer

import (
	"context"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/common/logger"
	"github.com/ozzy-cox/automatic-message-system/internal/common/queue"
	"github.com/redis/go-redis/v9"
)

const limit = 2
const offsetKey = "producer_offset"

type Service struct {
	Config            *ProducerConfig
	ProducerOnStatus  *atomic.Bool
	Cache             *redis.Client
	MessageRepository db.IMessageRepository
	Queue             *queue.WriterClient
	Logger            *logger.Logger
}

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

func (s *Service) ProduceMessages(wg *sync.WaitGroup, ctx context.Context, ticker *time.Ticker) {
	offset := s.mustGetProducerOffset()
	poffset := &offset

	for {
		select {
		case <-ctx.Done():
			s.Logger.Println("Context canceled, saving final offset")
			s.mustSetProducerOffset(poffset)
			wg.Done()
			return
		case <-ticker.C:
			if !s.ProducerOnStatus.Load() {
				continue
			}
			s.Logger.Printf("Fetching messages starting at offset: %d", offset)
			messages := s.MessageRepository.GetUnsentMessagesFromDb(2, offset)

			for msg, err := range messages {
				if err != nil {
					s.Logger.Printf("Error scanning messages: %v", err)
					continue
				}
				payload := queue.MessagePayload{
					ID:        msg.ID,
					Content:   msg.Content,
					To:        msg.To,
					CreatedAt: msg.CreatedAt,
				}

				if err := s.Queue.WriteMessage(ctx, payload); err != nil {
					s.Logger.Printf("Error writing message to queue: %v", err)
					continue
				}
				s.Logger.Printf("Successfully queued message ID: %d for recipient: %s", payload.ID, payload.To)
				(*poffset)++
			}
		}
	}
}

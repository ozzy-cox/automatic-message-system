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

func (s *Service) PushMessagesToQ(ctx context.Context, limit, offset int) int {
	s.Logger.Printf("Fetching messages starting at offset: %d", offset)
	messages := s.MessageRepository.GetUnsentMessagesFromDb(limit, offset)

	parsedMessages := make([]queue.MessagePayload, 0)
	for msg, err := range messages {
		if err != nil {
			s.Logger.Printf("Error scanning messages: %v", err)
			continue
		}
		parsedMessages = append(parsedMessages, queue.MessagePayload{
			ID:        msg.ID,
			Content:   msg.Content,
			To:        msg.To,
			CreatedAt: msg.CreatedAt,
		})
		s.Logger.Printf("Successfully queued message ID: %d for recipient: %s", msg.ID, msg.To)
	}

	if err := s.Queue.WriteMessages(ctx, parsedMessages...); err != nil {
		s.Logger.Printf("Error writing messages to queue: %v", err)
	}
	return len(parsedMessages)

}

func (s *Service) ProduceMessages(ctx context.Context, wg *sync.WaitGroup) {
	s.ProducerOnStatus.Store(true)
	ticker := time.NewTicker(s.Config.Interval)
	defer ticker.Stop()

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
			limit := s.Config.BatchCount
			go s.PushMessagesToQ(ctx, limit, offset)
			(*poffset) += limit
		}
	}
}

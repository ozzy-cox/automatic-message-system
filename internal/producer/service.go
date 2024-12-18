package producer

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
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
}

func (service *Service) mustGetProducerOffset() int {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	offset, err := service.Cache.Get(ctx, offsetKey).Result()
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

func (service *Service) mustSetProducerOffset(offsetValue *int) {
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Second)
	defer cancel()

	_, err := service.Cache.Set(ctx, offsetKey, *offsetValue, redis.KeepTTL).Result()
	if err != nil {
		panic("Can't set producer offset to redis")
	}
}

func (service *Service) ProduceMessages(wg *sync.WaitGroup, ctx context.Context, ticker *time.Ticker) {
	offset := service.mustGetProducerOffset()
	poffset := &offset
	for {
		select {
		case <-ctx.Done():
			service.mustSetProducerOffset(poffset)
			wg.Done()
			return
		case <-ticker.C:
			if !service.ProducerOnStatus.Load() {
				continue
			}
			fmt.Println("Producing", offset)
			messages := service.MessageRepository.GetUnsentMessagesFromDb(2, offset)

			for msg, err := range messages {
				msg := queue.MessagePayload{
					ID:        msg.ID,
					Content:   msg.Content,
					To:        msg.To,
					CreatedAt: msg.CreatedAt,
				}
				if err != nil {
					panic("Failed to scan messages")
				}

				service.Queue.WriteMessage(ctx, msg)
				(*poffset)++
			}
		}
	}
}

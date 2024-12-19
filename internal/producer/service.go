package producer

import (
	"sync/atomic"

	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/common/logger"
	"github.com/ozzy-cox/automatic-message-system/internal/common/queue"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	Config            *ProducerConfig
	ProducerOnStatus  *atomic.Bool
	Cache             *redis.Client
	MessageRepository db.MessageRepository
	Queue             *queue.WriterClient
	Logger            *logger.Logger
}

func NewProducerService(
	config *ProducerConfig,
	cache *redis.Client,
	messageRepository db.MessageRepository,
	queue *queue.WriterClient,
	logger *logger.Logger,
) *Service {
	return &Service{
		Config:            config,
		ProducerOnStatus:  &atomic.Bool{},
		Cache:             cache,
		MessageRepository: messageRepository,
		Queue:             queue,
		Logger:            logger,
	}
}

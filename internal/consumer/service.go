package consumer

import (
	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/common/logger"
	"github.com/ozzy-cox/automatic-message-system/internal/common/queue"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	Config            *ConsumerConfig
	Cache             *redis.Client
	MessageRepository db.MessageRepository
	ReaderQClient     queue.ReaderClient
	WriterQClient     queue.WriterClient
	RetryQueueWriter  queue.WriterClient
	Logger            *logger.Logger
}

func NewConsumerService(
	config *ConsumerConfig,
	cache *redis.Client,
	messageRepository db.MessageRepository,
	readerQClient queue.ReaderClient,
	writerQClient queue.WriterClient,
	retryQueueWriter queue.WriterClient,
	logger *logger.Logger,
) *Service {
	return &Service{
		Config:            config,
		Cache:             cache,
		MessageRepository: messageRepository,
		ReaderQClient:     readerQClient,
		WriterQClient:     writerQClient,
		RetryQueueWriter:  retryQueueWriter,
		Logger:            logger,
	}
}

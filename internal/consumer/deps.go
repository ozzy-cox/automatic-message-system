package consumer

import (
	"database/sql"
	"log"

	"github.com/ozzy-cox/automatic-message-system/internal/common/cache"
	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/common/logger"
	"github.com/ozzy-cox/automatic-message-system/internal/common/queue"
	"github.com/redis/go-redis/v9"
)

type ConsumerDeps struct {
	DBConnection           *sql.DB
	CacheClient            *redis.Client
	Logger                 *logger.Logger
	QueueWriterClient      queue.WriterClient
	QueueReaderClient      queue.ReaderClient
	RetryQueueWriterClient queue.WriterClient
}

func NewConsumerDeps(cfg ConsumerConfig) *ConsumerDeps {
	loggerInst, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	dbConnection, err := db.NewConnection(cfg.Database)
	if err != nil {
		loggerInst.Fatalf("Could not connect to db: %v", err)
	}

	cacheClient, err := cache.NewClient(cfg.Cache)
	if err != nil {
		loggerInst.Fatalf("Could not connect to cache: %v", err)
	}

	queueReaderClient, err := queue.NewKafkaReaderClient(cfg.Queue)
	if err != nil {
		loggerInst.Fatalf("Could not connect to queue: %v", err)
	}
	queueWriterClient, err := queue.NewKafkaWriterClient(cfg.Queue)
	if err != nil {
		loggerInst.Fatalf("Could not connect to queue: %v", err)
	}
	retryQueueWriterClient, err := queue.NewKafkaWriterClient(cfg.RetryQueue)
	if err != nil {
		loggerInst.Fatalf("Could not connect to retry-queue: %v", err)
	}

	return &ConsumerDeps{
		DBConnection:           dbConnection,
		CacheClient:            cacheClient,
		Logger:                 loggerInst,
		QueueWriterClient:      queueWriterClient,
		QueueReaderClient:      queueReaderClient,
		RetryQueueWriterClient: retryQueueWriterClient,
	}
}

func (d *ConsumerDeps) Cleanup() {
	d.QueueWriterClient.Close()
	d.QueueReaderClient.Close()
	d.RetryQueueWriterClient.Close()
}

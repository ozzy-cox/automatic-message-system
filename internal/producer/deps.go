package producer

import (
	"database/sql"
	"log"

	"github.com/ozzy-cox/automatic-message-system/internal/common/cache"
	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/common/logger"
	"github.com/ozzy-cox/automatic-message-system/internal/common/queue"
	"github.com/redis/go-redis/v9"
)

type ProducerDeps struct {
	DBConnection      *sql.DB
	CacheClient       *redis.Client
	QueueWriterClient queue.WriterClient
	Logger            *logger.Logger
}

func NewProducerDeps(cfg ProducerConfig) *ProducerDeps {
	loggerInst, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	dbConnection, err := db.NewConnection(cfg.Database)
	if err != nil {
		loggerInst.Fatalf("Could not load config: %v", err)
	}

	cacheClient, err := cache.NewClient(cfg.Cache)
	if err != nil {
		loggerInst.Fatalf("Could not connect to cache: %v", err)
	}

	queueClient, err := queue.NewKafkaWriterClient(cfg.Queue)
	if err != nil {
		loggerInst.Fatalf("Could not connect to queue: %v", err)
	}

	return &ProducerDeps{
		DBConnection:      dbConnection,
		CacheClient:       cacheClient,
		QueueWriterClient: queueClient,
		Logger:            loggerInst,
	}
}

func (d *ProducerDeps) Cleanup() {
	d.QueueWriterClient.Close()
	d.Logger.Println("Cleaning up")
}

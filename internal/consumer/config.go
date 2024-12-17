package consumer

import (
	"time"

	"github.com/ozzy-cox/automatic-message-system/internal/common/cache"
	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/common/queue"
	"github.com/ozzy-cox/automatic-message-system/internal/common/utils"
)

type ConsumerConfig struct {
	Database   db.DatabaseConfig
	Cache      cache.RedisConfig
	Queue      queue.KafkaConfig
	Interval   time.Duration
	Port       string
	RequestURL string
}

func GetConsumerConfig() (*ConsumerConfig, error) {
	config := &ConsumerConfig{
		RequestURL: utils.GetEnvStringWithDefault("REQUEST_URL", "http://localhost:3000"),
		Database: db.DatabaseConfig{
			Host:     utils.GetEnvStringWithDefault("DB_HOST", "localhost"),
			Port:     utils.GetEnvStringWithDefault("DB_PORT", "5432"),
			User:     utils.GetEnvStringWithDefault("DB_USER", "postgres"),
			Password: utils.GetEnvStringWithDefault("DB_PASSWORD", "postgres"),
			DBName:   utils.GetEnvStringWithDefault("DB_NAME", "automatic_message_system"),
			SSLMode:  utils.GetEnvStringWithDefault("DB_SSLMODE", "disable"),
		},
		Cache: cache.RedisConfig{
			Host: utils.GetEnvStringWithDefault("REDIS_HOST", "localhost"),
			Port: utils.GetEnvStringWithDefault("REDIS_PORT", "6379"),
			DB:   utils.GetEnvIntWithDefault("REDIS_DB", 0),
		},
		Queue: queue.KafkaConfig{
			Brokers: queue.GetDefaultKafkaBrokers(),
			GroupID: utils.GetEnvStringWithDefault("KAFKA_GROUP_ID", "message-consumer"),
			Topic:   utils.GetEnvStringWithDefault("KAFKA_TOPIC", "messages"),
		},
	}

	return config, nil
}

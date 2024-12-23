package producer

import (
	"time"

	"github.com/ozzy-cox/automatic-message-system/internal/common/cache"
	"github.com/ozzy-cox/automatic-message-system/internal/common/config"
	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/common/logger"
	"github.com/ozzy-cox/automatic-message-system/internal/common/queue"
	"github.com/ozzy-cox/automatic-message-system/internal/common/utils"
)

type ProducerConfig struct {
	config.BaseConfig
	Queue      queue.KafkaConfig
	Interval   time.Duration
	BatchCount int
	Port       string
}

func GetProducerConfig() (*ProducerConfig, error) {
	config := &ProducerConfig{
		BaseConfig: config.BaseConfig{
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
			Logger: logger.Config{
				LogFile:     utils.GetEnvStringWithDefault("LOG_FILE", "/tmp/log/automatic-message-system/producer.log"),
				LogToStdout: utils.GetEnvBoolWithDefault("LOG_TO_STDOUT", true),
			},
		},
		Queue: queue.KafkaConfig{
			Brokers: queue.GetDefaultKafkaBrokers(),
			GroupID: utils.GetEnvStringWithDefault("KAFKA_GROUP_ID", "message-consumer"),
			Topic:   utils.GetEnvStringWithDefault("KAFKA_TOPIC", "messages"),
		},
		Interval:   utils.GetEnvDurationWithDefault("INTERVAL", time.Minute*2),
		BatchCount: utils.GetEnvIntWithDefault("BATCH_COUNT", 1),
		Port:       utils.GetEnvStringWithDefault("PRODUCER_API_PORT", "8001"),
	}

	return config, nil
}

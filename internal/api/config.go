package api

import (
	"github.com/ozzy-cox/automatic-message-system/internal/common/cache"
	"github.com/ozzy-cox/automatic-message-system/internal/common/config"
	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/common/logger"
	"github.com/ozzy-cox/automatic-message-system/internal/common/utils"
)

type APIConfig struct {
	config.BaseConfig
	ProducerURL string
	Host        string
	Port        string
}

func GetAPIConfig() (*APIConfig, error) {
	config := &APIConfig{
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
				LogFile:     utils.GetEnvStringWithDefault("LOG_FILE", "/tmp/log/automatic-message-system/api.log"),
				LogToStdout: utils.GetEnvBoolWithDefault("LOG_TO_STDOUT", true),
			},
		},
		ProducerURL: utils.GetEnvStringWithDefault("PRODUCER_URL", "http://localhost:8001"),
		Host:        utils.GetEnvStringWithDefault("HOST", "127.0.0.1"),
		Port:        utils.GetEnvStringWithDefault("PORT", "8080"),
	}

	return config, nil
}

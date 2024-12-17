package config

import (
	"time"
)

var WorkerConfigObject *WorkerConfig

type WorkerConfig struct {
	Database DatabaseConfig
	Cache    RedisConfig
	Interval time.Duration
	Port     string
}

func GetWorkerConfig() (*WorkerConfig, error) {
	if WorkerConfigObject != nil {
		return WorkerConfigObject, nil
	}

	config := &WorkerConfig{
		Interval: getEnvDurationWithDefault("WORKER_INTERVAL", time.Second),
		Port:     getEnvStringWithDefault("WORKER_PORT", "8001"),
		Database: DatabaseConfig{
			Host:     getEnvStringWithDefault("DB_HOST", "localhost"),
			Port:     getEnvStringWithDefault("DB_PORT", "5432"),
			User:     getEnvStringWithDefault("DB_USER", "postgres"),
			Password: getEnvStringWithDefault("DB_PASSWORD", "postgres"),
			DBName:   getEnvStringWithDefault("DB_NAME", "automatic_message_system"),
			SSLMode:  getEnvStringWithDefault("DB_SSLMODE", "disable"),
		},
		Cache: RedisConfig{
			Host: getEnvStringWithDefault("REDIS_HOST", "localhost"),
			Port: getEnvStringWithDefault("REDIS_PORT", "6379"),
			DB:   getEnvIntWithDefault("REDIS_DB", 0),
		}}

	WorkerConfigObject = config
	return config, nil
}

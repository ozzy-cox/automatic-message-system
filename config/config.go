package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTP     HTTPConfig
	Worker   WorkerConfig
	Database DatabaseConfig
	Cache    RedisConfig
}

type HTTPConfig struct {
	Host string
	Port string
}

type WorkerConfig struct {
	Interval time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host string
	Port string
	DB   int
}

func Load() (*Config, error) {
	config := &Config{
		HTTP: HTTPConfig{
			Host: getEnvStringWithDefault("HOST", "127.0.0.1"),
			Port: getEnvStringWithDefault("PORT", "8080"),
		},
		Worker: WorkerConfig{
			Interval: getEnvDurationWithDefault("WORKER_INTERVAL", time.Second),
		},
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

	return config, nil
}

func getEnvStringWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvDurationWithDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsedValue, err := strconv.Atoi(value); err == nil {
			return parsedValue
		}
	}
	return defaultValue
}

package config

import (
	"os"
	"time"
)

type Config struct {
	HTTP     HTTPConfig
	Worker   WorkerConfig
	Database DatabaseConfig
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

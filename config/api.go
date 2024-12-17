package config

type APIConfig struct {
	HTTP      HTTPConfig
	Database  DatabaseConfig
	Cache     RedisConfig
	WorkerUrl string
}

type HTTPConfig struct {
	Host string
	Port string
}

var APIConfigObject *APIConfig

func GetAPIConfig() (*APIConfig, error) {
	if APIConfigObject != nil {
		return APIConfigObject, nil
	}
	config := &APIConfig{
		HTTP: HTTPConfig{
			Host: getEnvStringWithDefault("HOST", "127.0.0.1"),
			Port: getEnvStringWithDefault("PORT", "8080"),
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
		},
		WorkerUrl: getEnvStringWithDefault("WORKER_URL", "http://localhost:8001"),
	}

	APIConfigObject = config

	return config, nil
}

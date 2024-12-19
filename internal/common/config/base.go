package config

import (
	"github.com/ozzy-cox/automatic-message-system/internal/common/cache"
	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/common/logger"
)

type BaseConfig struct {
	Database db.DatabaseConfig
	Cache    cache.RedisConfig
	Logger   logger.Config
}

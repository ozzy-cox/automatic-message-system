package api

import (
	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/common/logger"
)

type Service struct {
	Config            *APIConfig
	MessageRepository db.MessageRepository
	Logger            *logger.Logger
}

func NewAPIService(
	config *APIConfig,
	messageRepository db.MessageRepository,
	logger *logger.Logger,
) *Service {
	return &Service{
		Config:            config,
		MessageRepository: messageRepository,
		Logger:            logger,
	}
}

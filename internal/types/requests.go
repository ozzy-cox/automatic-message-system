package types

import "github.com/ozzy-cox/automatic-message-system/internal/models"

type ToggleRequest struct {
	// NOTE This should be the desired status of the worker
	WorkerStatus *bool `json:"workerStatus"`
}

type ToggleResponse struct {
	// NOTE This should be the current status of the worker
	WorkerStatus bool `json:"WorkerStatus"`
}

type SentMessagesResponse struct {
	SentMessages []models.Message `json:"sentMessages"`
}

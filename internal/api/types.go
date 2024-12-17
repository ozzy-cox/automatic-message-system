package api

import "github.com/ozzy-cox/automatic-message-system/internal/common/db"

type ToggleRequest struct {
	// NOTE This should be the desired status of the worker
	WorkerStatus *bool `json:"workerStatus"`
}

type ToggleResponse struct {
	// NOTE This should be the current status of the worker
	WorkerStatus bool `json:"WorkerStatus"`
}

type SentMessagesResponse struct {
	// TODO change the object in the response, dont use db types
	SentMessages []db.Message `json:"sentMessages"`
}

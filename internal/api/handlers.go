package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/common/logger"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	Config            *APIConfig
	MessageRepository db.IMessageRepository
	Logger            *logger.Logger
}

// HandleGetSentMessages godoc
//
//	@Summary		Get sent messages
//	@Description	Retrieves a list of sent messages from the system
//	@Produce		json
//	@Success		200	{object}	SentMessagesResponse
//	@Router			/sent-messages [get]
func (s *Service) HandleGetSentMessages(w http.ResponseWriter, r *http.Request) {
	// TODO Add pagination
	rows := s.MessageRepository.GetSentMessagesFromDb(20, 0)

	sentMessages := make([]SentMessage, 0)
	for i, err := range rows {
		if err != nil {
			http.Error(w, "Failed to scan messages", http.StatusInternalServerError)
			return
		}
		sentMessages = append(sentMessages, SentMessage(*i))
	}

	response := SentMessagesResponse{
		SentMessages: sentMessages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleToggleWorker godoc
//
//	@Summary		Toggle worker status
//	@Description	Toggles the message sending worker on/off
//	@Tags			worker
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ToggleRequest	true	"Worker status toggle request"
//	@Success		200		{object}	ToggleResponse
//	@Failure		400		{string}	string	"Invalid request body"
//	@Router			/toggle-worker [post]
func (s *Service) HandleToggleWorker(w http.ResponseWriter, r *http.Request) {
	var request ToggleRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.WorkerStatus == nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	workerUrl := s.Config.ProducerURL + "/toggle-worker"

	jsonBody, _ := json.Marshal(request)

	http.Post(workerUrl, "application/json", bytes.NewReader(jsonBody))
}

package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/ozzy-cox/automatic-message-system/internal/common/utils"
)

// HandleGetSentMessages godoc
//
// @Summary      Get sent messages
// @Description  Retrieves a list of sent messages from the system
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        limit  query     integer  false  "Number of messages to return (default: 20)"  minimum(1)  maximum(100)
// @Param        offset query     integer  false  "Number of messages to skip (default: 0)"     minimum(0)
// @Success      200    {object}  utils.JsonResponse{data=SentMessagesResponse}
// @Failure      500    {object}  utils.JsonResponse{message=string}
// @Router       /sent-messages [get]
func (s *Service) HandleGetSentMessages(w http.ResponseWriter, r *http.Request) {
	params := GetSentMessagesParams{
		Limit:  utils.GetIntParam(r, "limit", 20),
		Offset: utils.GetIntParam(r, "offset", 0),
	}
	rows := s.MessageRepository.GetSentMessages(params.Limit, params.Offset)

	sentMessages := make([]SentMessage, 0)
	for i, err := range rows {
		if err != nil {
			http.Error(w, "Failed to scan messages", http.StatusInternalServerError)
			return
		}
		sentMessages = append(sentMessages, SentMessage(*i))
	}

	response := utils.JsonResponse{
		Data: SentMessagesResponse{
			SentMessages: sentMessages,
		},
		Message: "success",
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

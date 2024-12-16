package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/ozzy-cox/automatic-message-system/internal/models"
	"github.com/ozzy-cox/automatic-message-system/internal/types"
)

var isWorkerRunning atomic.Bool

func HandleGetSentMessages(w http.ResponseWriter, r *http.Request) {
	sentMessages := []models.Message{
		{ID: 1, Content: "asdf", To: "+901233453434", Sent: true, SentAt: time.Now(), CreatedAt: time.Now()},
		{ID: 2, Content: "asdf", To: "+901233453434", Sent: true, SentAt: time.Now(), CreatedAt: time.Now()},
		{ID: 3, Content: "asdf", To: "+901233453434", Sent: true, SentAt: time.Now(), CreatedAt: time.Now()},
		{ID: 4, Content: "asdf", To: "+901233453434", Sent: true, SentAt: time.Now(), CreatedAt: time.Now()},
		{ID: 5, Content: "asdf", To: "+901233453434", Sent: true, SentAt: time.Now(), CreatedAt: time.Now()},
	}

	response := types.SentMessagesResponse{
		SentMessages: sentMessages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HandleToggleWorker(w http.ResponseWriter, r *http.Request) {
	var request types.ToggleRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.WorkerStatus == nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("request", request)

	isWorkerRunning.Store(*request.WorkerStatus)
	fmt.Println(isWorkerRunning.Load())

	fmt.Println("Handle toggle worker")
}

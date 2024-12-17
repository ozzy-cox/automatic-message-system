package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ozzy-cox/automatic-message-system/config"
	"github.com/ozzy-cox/automatic-message-system/internal/db"
	"github.com/ozzy-cox/automatic-message-system/internal/types"
)

func HandleGetSentMessages(w http.ResponseWriter, r *http.Request) {
	// TODO Add pagination
	rows, err := db.DbConnection.Query("SELECT * FROM messages LIMIT 20")
	if err != nil {
		fmt.Println("Error getting messages from db")
	}

	sentMessages := make([]db.Message, 0)
	for rows.Next() {
		var msg db.Message
		err := rows.Scan(
			&msg.ID,
			&msg.Content,
			&msg.To,
			&msg.Sent,
			&msg.SentAt,
			&msg.CreatedAt,
		)
		if err != nil {
			http.Error(w, "Failed to scan messages", http.StatusInternalServerError)
			return
		}
		sentMessages = append(sentMessages, msg)
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

	workerUrl := config.APIConfigObject.WorkerUrl + "/toggle-worker"

	jsonBody, _ := json.Marshal(request)

	http.Post(workerUrl, "application/json", bytes.NewReader(jsonBody))

}

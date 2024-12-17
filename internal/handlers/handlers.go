package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/ozzy-cox/automatic-message-system/internal/db"
	"github.com/ozzy-cox/automatic-message-system/internal/types"
)

var (
	dbConn          *sql.DB
	isWorkerRunning atomic.Bool
)

func Initialize(_dbConn *sql.DB) {
	dbConn = _dbConn
}

func HandleGetSentMessages(w http.ResponseWriter, r *http.Request) {
	// TODO Add pagination
	rows, err := dbConn.Query("SELECT * FROM messages LIMIT 20")
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

	isWorkerRunning.Store(*request.WorkerStatus)
}

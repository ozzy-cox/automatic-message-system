package producer

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sync/atomic"

	"github.com/ozzy-cox/automatic-message-system/internal/common/queue"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	Config           *ProducerConfig
	ProducerOnStatus *atomic.Bool
	Cache            *redis.Client
	DB               *sql.DB
	Queue            *queue.WriterClient
}

func (app *Service) HandleToggleProducer(w http.ResponseWriter, r *http.Request) {
	var request toggleRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Bad request body", http.StatusBadRequest)
	}

	app.ProducerOnStatus.Store(request.WorkerStatus)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r.Body)
}

package worker

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
)

var ProducerOnStatus atomic.Bool

type toggleRequest struct {
	WorkerStatus bool `json:"workerStatus"`
}

func HandleToggleProducer(w http.ResponseWriter, r *http.Request) {
	var request toggleRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Bad request body", http.StatusBadRequest)
	}

	ProducerOnStatus.Store(request.WorkerStatus)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r.Body)
}

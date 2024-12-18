package producer

import (
	"encoding/json"
	"net/http"
)

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

package main

import (
	"github.com/ozzy-cox/automatic-message-system/config"
	"github.com/ozzy-cox/automatic-message-system/internal/handlers"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	http.HandleFunc("GET /sent-messages", handlers.HandleGetSentMessages)
	http.HandleFunc("POST /toggle-worker", handlers.HandleToggleWorker)

	addr := ":" + cfg.HTTP.Port
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

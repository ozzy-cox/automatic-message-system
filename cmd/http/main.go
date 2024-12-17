package main

import (
	"log"
	"net/http"

	"github.com/ozzy-cox/automatic-message-system/config"
	"github.com/ozzy-cox/automatic-message-system/internal/db"
	"github.com/ozzy-cox/automatic-message-system/internal/handlers"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	dbConn, err := db.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Could not load database: %v", err)
	}

	handlers.Initialize(dbConn)

	http.HandleFunc("GET /sent-messages", handlers.HandleGetSentMessages)
	http.HandleFunc("POST /toggle-worker", handlers.HandleToggleWorker)

	addr := ":" + cfg.HTTP.Port
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

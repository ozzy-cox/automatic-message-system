package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/producer"
)

func main() {
	cfg, err := producer.GetProducerConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	deps := producer.NewProducerDeps(*cfg)
	defer deps.Cleanup()

	service := producer.NewProducerService(
		cfg,
		deps.CacheClient,
		db.NewMessageRepository(deps.DBConnection),
		deps.QueueWriterClient,
		deps.Logger,
	)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)

	go service.ProduceMessages(ctx, &wg)

	http.HandleFunc("POST /toggle-worker", service.HandleToggleProducer)
	addr := ":" + cfg.Port
	go func() {
		deps.Logger.Printf("Starting HTTP server on %s", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			deps.Logger.Fatalf("HTTP server error: %v", err)
		}
	}()

	deps.Logger.Println("Producer service started")
	<-stop
	deps.Logger.Println("Shutting down...")
	cancel()
	wg.Wait()
}

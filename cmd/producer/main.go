package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ozzy-cox/automatic-message-system/internal/common/cache"
	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/common/logger"
	"github.com/ozzy-cox/automatic-message-system/internal/common/queue"
	"github.com/ozzy-cox/automatic-message-system/internal/producer"
)

func main() {
	cfg, err := producer.GetProducerConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	loggerInst, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	dbConn, err := db.NewConnection(cfg.Database)
	if err != nil {
		loggerInst.Fatalf("Could not load config: %v", err)
	}

	cacheClient, err := cache.NewClient(cfg.Cache)
	if err != nil {
		loggerInst.Fatalf("Could not connect to cache: %v", err)
	}

	queueClient, err := queue.NewWriterClient(cfg.Queue)
	if err != nil {
		loggerInst.Fatalf("Could not connect to queue: %v", err)
	}
	defer queueClient.Close()

	service := producer.NewProducerService(
		cfg,
		cacheClient,
		db.NewMessageRepository(dbConn),
		queueClient,
		loggerInst,
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
		loggerInst.Printf("Starting HTTP server on %s", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			loggerInst.Fatalf("HTTP server error: %v", err)
		}
	}()

	loggerInst.Println("Producer service started")
	<-stop
	loggerInst.Println("Shutting down...")
	cancel()
	wg.Wait()
}

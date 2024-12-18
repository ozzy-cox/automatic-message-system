package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

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

	service := producer.Service{
		Config:           cfg,
		ProducerOnStatus: &atomic.Bool{},
		Cache:            cacheClient,
		MessageRepository: &db.MessageRepository{
			DB: dbConn,
		},
		Queue:  queueClient,
		Logger: loggerInst,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	ticker := time.NewTicker(cfg.Interval.Abs())
	defer ticker.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)

	service.ProducerOnStatus.Store(true)
	go service.ProduceMessages(&wg, ctx, ticker)

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
	if err := queueClient.Close(); err != nil {
		loggerInst.Printf("Error closing queue client: %v", err)
	}
}

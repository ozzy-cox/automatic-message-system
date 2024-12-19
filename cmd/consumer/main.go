package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ozzy-cox/automatic-message-system/internal/common/cache"
	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/common/logger"
	"github.com/ozzy-cox/automatic-message-system/internal/common/queue"
	"github.com/ozzy-cox/automatic-message-system/internal/consumer"
)

func main() {
	cfg, err := consumer.GetConsumerConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	loggerInst, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	dbConn, err := db.NewConnection(cfg.Database)
	if err != nil {
		loggerInst.Fatalf("Could not connect to db: %v", err)
	}

	cacheClient, err := cache.NewClient(cfg.Cache)
	if err != nil {
		loggerInst.Fatalf("Could not connect to cache: %v", err)
	}

	readQueueClient, err := queue.NewReaderClient(cfg.Queue)
	if err != nil {
		loggerInst.Fatalf("Could not connect to queue: %v", err)
	}
	defer readQueueClient.Close()
	writeQueueClient, err := queue.NewWriterClient(cfg.Queue)
	if err != nil {
		loggerInst.Fatalf("Could not connect to queue: %v", err)
	}
	defer writeQueueClient.Close()
	retryQueueWriterClient, err := queue.NewWriterClient(cfg.RetryQueue)
	if err != nil {
		loggerInst.Fatalf("Could not connect to retry-queue: %v", err)
	}
	defer retryQueueWriterClient.Close()

	service := consumer.Service{
		Config:            cfg,
		Cache:             cacheClient,
		MessageRepository: db.NewMessageRepository(dbConn),
		ReaderQClient:     readQueueClient,
		WriterQClient:     writeQueueClient,
		RetryQueueWriter:  retryQueueWriterClient,
		Logger:            loggerInst,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	go service.ConsumeMessages(ctx, &wg)

	loggerInst.Println("Consumer service started")
	<-stop
	loggerInst.Println("Shutting down ...")
	cancel()
	wg.Wait()
}

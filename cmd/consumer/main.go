package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/consumer"
)

func main() {
	cfg, err := consumer.GetConsumerConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}
	deps := consumer.NewConsumerDeps(*cfg)
	defer deps.Cleanup()

	service := consumer.NewConsumerService(
		cfg,
		deps.CacheClient,
		db.NewMessageRepository(deps.DBConnection),
		deps.QueueReaderClient,
		deps.QueueWriterClient,
		deps.RetryQueueWriterClient,
		deps.Logger,
	)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	go service.ConsumeMessages(ctx, &wg)

	deps.Logger.Println("Consumer service started")
	<-stop
	deps.Logger.Println("Shutting down ...")
	cancel()
	wg.Wait()
}

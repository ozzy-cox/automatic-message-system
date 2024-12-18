package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ozzy-cox/automatic-message-system/internal/common/cache"
	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/common/queue"
	"github.com/ozzy-cox/automatic-message-system/internal/consumer"
)

func main() {
	cfg, err := consumer.GetConsumerConfig()
	if err != nil {
		fmt.Printf("Could not load config: %v\n", err)
	}

	dbConn, err := db.NewConnection(cfg.Database)
	if err != nil {
		fmt.Printf("Could not connect to db: %v\n", err)
		panic(err)
	}

	cacheClient, err := cache.NewClient(cfg.Cache)
	if err != nil {
		fmt.Printf("Could not connect to cache: %v\n", err)
		panic(err)
	}

	queueClient, err := queue.NewReaderClient(cfg.Queue)
	if err != nil {
		fmt.Printf("Could not connect to cache: %v\n", err)
		panic(err)
	}
	defer queueClient.Close()

	service := consumer.Service{
		Config: cfg,
		DB:     dbConn,
		Cache:  cacheClient,
		Queue:  queueClient,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	go service.ConsumeMessages(ctx, &wg)
	<-stop
	cancel()
	wg.Wait()
}

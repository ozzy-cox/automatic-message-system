package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ozzy-cox/automatic-message-system/config"
	"github.com/ozzy-cox/automatic-message-system/internal/cache"
	"github.com/ozzy-cox/automatic-message-system/internal/db"
	"github.com/ozzy-cox/automatic-message-system/internal/handlers"
	"github.com/ozzy-cox/automatic-message-system/internal/worker"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Could not load config: %v\n", err)
	}

	dbConn, err := db.NewConnection(cfg.Database)
	if err != nil {
		fmt.Printf("Could not connect to db: %v\n", err)
		panic(err)
	}
	handlers.Initialize(dbConn)

	cacheClient, err := cache.NewClient(cfg.Cache)
	if err != nil {
		fmt.Printf("Could not connect to cache: %v\n", err)
		panic(err)
	}
	worker.Initialize(dbConn, cacheClient)
	messageChan := make(chan db.Message, 1000) // FIXME: This is simulating an interprocess queue rn

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(cfg.Worker.Interval.Abs())
	defer ticker.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(2)

	go worker.ProduceMessages(&wg, ctx, messageChan, ticker)
	go worker.ConsumeMessages(&wg, ctx, messageChan)

	<-stop
	cancel()

	wg.Wait()
}

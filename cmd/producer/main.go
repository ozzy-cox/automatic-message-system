package main

import (
	"context"
	"fmt"
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
	"github.com/ozzy-cox/automatic-message-system/internal/common/queue"
	"github.com/ozzy-cox/automatic-message-system/internal/producer"
)

func main() {
	cfg, err := producer.GetProducerConfig()
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

	queueClient, err := queue.NewWriterClient(cfg.Queue)
	if err != nil {
		fmt.Printf("Could not connect to cache: %v\n", err)
		panic(err)
	}

	service := producer.Service{
		Config:           cfg,
		ProducerOnStatus: &atomic.Bool{},
		Cache:            cacheClient,
		DB:               dbConn,
		Queue:            queueClient,
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
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	<-stop
	cancel()
	fmt.Println("Exiting gracefully...")
	wg.Wait()
	queueClient.Close()
}

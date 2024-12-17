package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ozzy-cox/automatic-message-system/config"
	"github.com/ozzy-cox/automatic-message-system/internal/cache"
	"github.com/ozzy-cox/automatic-message-system/internal/db"
	"github.com/ozzy-cox/automatic-message-system/internal/worker"
)

func main() {
	cfg, err := config.GetWorkerConfig()
	if err != nil {
		fmt.Printf("Could not load config: %v\n", err)
	}

	_, err = db.GetConnection(cfg.Database)
	if err != nil {
		fmt.Printf("Could not connect to db: %v\n", err)
		panic(err)
	}

	_, err = cache.GetClient(cfg.Cache)
	if err != nil {
		fmt.Printf("Could not connect to cache: %v\n", err)
		panic(err)
	}
	messageChan := make(chan db.Message, 1000) // FIXME: This is simulating an interprocess queue rn

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(cfg.Interval.Abs())
	defer ticker.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(2)

	go worker.ProduceMessages(&wg, ctx, messageChan, ticker)
	go worker.ConsumeMessages(&wg, ctx, messageChan)
	worker.ProducerOnStatus.Store(true)

	http.HandleFunc("POST /toggle-worker", worker.HandleToggleProducer)

	addr := ":" + cfg.Port
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	<-stop
	cancel()

	wg.Wait()
}

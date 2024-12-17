package main

import (
	"fmt"
	"time"

	"github.com/ozzy-cox/automatic-message-system/config"
	"github.com/ozzy-cox/automatic-message-system/internal/db"
	"github.com/ozzy-cox/automatic-message-system/internal/handlers"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Could not load config: %v\n", err)
	}

	dbConn, err := db.NewConnection(cfg.Database)
	if err != nil {
		fmt.Printf("Could not connect to db: %v\n", err)
	}
	handlers.Initialize(dbConn)

	ticker := time.NewTicker(cfg.Worker.Interval.Abs())
	defer ticker.Stop()

	for ; true; <-ticker.C {
		do()
	}

}

func do() {
	fmt.Println("hello world!")
}

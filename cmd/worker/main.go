package main

import (
	"fmt"
	"time"

	"github.com/ozzy-cox/automatic-message-system/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Could not load config: %v\n", err)
	}

	ticker := time.NewTicker(cfg.Worker.Interval.Abs())
	defer ticker.Stop()

	for ; true; <-ticker.C {
		do()
	}

}

func do() {
	fmt.Println("hello world!")
}

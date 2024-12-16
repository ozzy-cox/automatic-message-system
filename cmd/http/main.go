package main

import (
	"fmt"
	"github.com/ozzy-cox/automatic-message-system/config"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "helloworld")
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	http.HandleFunc("/", handler)
	addr := ":" + cfg.HTTP.Port
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

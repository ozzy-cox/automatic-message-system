package main

import (
	"encoding/json"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
)

func handler(w http.ResponseWriter, r *http.Request) {
	errorRate := float64(0)
	errorRateEnv := os.Getenv("ERROR_RATE")
	if errorRateEnv != "" {
		parsedValue, err := strconv.ParseFloat(errorRateEnv, 64)
		if err != nil {
			errorRate = 0
		}
		errorRate = parsedValue
	}
	if rand.Float64() < errorRate {
		log.Printf("returned random error")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Log the request details
	defer r.Body.Close()

	byt, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	var body map[string]any
	err = json.Unmarshal(byt, &body)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Received request: Method=%s, URL=%s, Host=%s, RemoteAddr=%s Body= %s",
		r.Method, r.URL.String(), r.Host, r.RemoteAddr, body)

	response := map[string]string{
		"message":   "Accepted",
		"messageId": uuid.NewString(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Set up a simple HTTP handler for all paths
	http.HandleFunc("/", handler)

	// Start listening on port 8080 and log any errors
	log.Println("Starting server on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}

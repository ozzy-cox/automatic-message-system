package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/ozzy-cox/automatic-message-system/internal/db"
)

// const requestURL = "https://webhook.site/5770a369-afb4-47dc-8d6b-d3da51530c81"
const requestURL = "http://localhost:3000"

func sendMessage(msg db.Message) {
	jsonBody, err := json.Marshal(msg)
	bodyReader := bytes.NewReader(jsonBody)

	resp, err := http.Post(requestURL, "application/json", bodyReader)
	if err != nil {
		// TODO requeue failed messages
		fmt.Println(err)
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		// TODO requeue failed messages
		fmt.Println("Error reading response:", err)
		return
	}

}

func ConsumeMessages(wg *sync.WaitGroup, ctx context.Context, messageChan chan db.Message) {
	for {
		select {
		case msg := <-messageChan:
			sendMessage(msg)
		case <-ctx.Done():
			wg.Done()
			return
		}
	}

}

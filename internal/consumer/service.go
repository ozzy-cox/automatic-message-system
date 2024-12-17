package consumer

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/common/queue"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	Config *ConsumerConfig
	Cache  *redis.Client
	DB     *sql.DB
	Queue  *queue.ReaderClient
}

func (service *Service) sendMessage(msg db.Message) {
	jsonBody, err := json.Marshal(msg)
	bodyReader := bytes.NewReader(jsonBody)

	resp, err := http.Post(service.Config.RequestURL, "application/json", bodyReader)
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

func (service *Service) ConsumeMessages(ctx context.Context, wg *sync.WaitGroup) {
	messageChan := make(chan db.Message, 1000)
	go func() {
		for {
			msg, err := service.Queue.ReadMessage(ctx)
			if err != nil {
				fmt.Println("Error reading message from kafka:", err)
				panic(err)
			}
			messageChan <- msg
		}
	}()
	for {
		select {
		case msg := <-messageChan:
			service.sendMessage(msg)
		case <-ctx.Done():
			wg.Done()
			return
		}
	}

}

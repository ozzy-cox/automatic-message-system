package consumer

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/ozzy-cox/automatic-message-system/internal/common/queue"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	Config *ConsumerConfig
	Cache  *redis.Client
	DB     *sql.DB
	Queue  *queue.ReaderClient
}

func (service *Service) sendMessage(msg queue.MessagePayload) {
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
	messageChan := make(chan queue.MessagePayload, 1000)
	go func() {
		for {
			msg, err := service.Queue.ReadMessage(ctx)
			if err != nil {
				close(messageChan)
				if errors.Is(err, context.Canceled) {
					return
				}
				fmt.Println("Error reading message from kafka:", err)
				panic(err)
			}
			messageChan <- msg
		}
	}()
	for {
		select {
		case msg, ok := <-messageChan:
			if !ok {
				wg.Done()
				return
			}
			service.sendMessage(msg)
		case <-ctx.Done():
			for {
				select {
				case msg, ok := <-messageChan:
					if !ok {
						wg.Done()
						return
					}
					service.sendMessage(msg)
				default:
					wg.Done()
					return
				}
			}
		}
	}

}

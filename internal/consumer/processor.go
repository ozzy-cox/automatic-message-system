package consumer

import (
	"context"
	"errors"
	"sync"

	"github.com/ozzy-cox/automatic-message-system/internal/common/queue"
)

func (s *Service) ReadMessages(ctx context.Context, messageChan chan queue.MessagePayload) {
	for {
		msg, err := s.ReaderQClient.ReadMessage(ctx)
		if err != nil {
			close(messageChan)
			if errors.Is(err, context.Canceled) {
				s.Logger.Println("Context canceled, stopping message reader")
				return
			}
			s.Logger.Printf("Error reading message from kafka: %v", err)
			panic(err)
		}
		messageChan <- msg
	}
}

func (s *Service) ConsumeMessages(ctx context.Context, wg *sync.WaitGroup) {
	messageChan := make(chan queue.MessagePayload, 1000)
	go s.ReadMessages(ctx, messageChan)
	for msg := range messageChan {
		go s.handleMessage(ctx, msg)
	}
	wg.Done()
}

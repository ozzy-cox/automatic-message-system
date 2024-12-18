package retryconsumer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"

	"github.com/ozzy-cox/automatic-message-system/internal/common/queue"
	"github.com/ozzy-cox/automatic-message-system/internal/consumer"
)

type Service struct {
	consumer.Service
	RetryQueueReader *queue.ReaderClient
	DLQueueWriter    *queue.WriterClient
}

func (s *Service) sendMessage(msg queue.MessagePayload) error { // TODO add context and cancellation

	jsonBody, err := json.Marshal(msg)
	if err != nil {
		s.Logger.Printf("Error marshaling message: %v", err)
		return err
	}
	bodyReader := bytes.NewReader(jsonBody)

	resp, err := http.Post(s.Config.RequestURL, "application/json", bodyReader)
	if err != nil || resp.StatusCode != 200 {
		s.Logger.Printf("Error sending message to %s: %v", s.Config.RequestURL, err)
		return err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		s.Logger.Printf("Error reading response body: %v", err)
		return nil
	}
	// s.Logger.Printf("Successfully sent message ID: %d to %s", msg.ID, msg.To)

	err = s.MessageRepository.SetMessageSent(msg.ID)
	if err != nil {
		s.Logger.Printf("Failed to update message sent state: %v", err)
		return nil
	}
	return nil
}

func (s *Service) repeatSendMessage(ctx context.Context, msg queue.MessagePayload) {
	for i := 1; i < s.Config.RetryCount; i++ {
		jitter := (rand.Float64() * 1) - 1
		nextWaitDuration := time.Duration(math.Pow(2, float64(i))+jitter) * s.Config.Interval
		err := s.sendMessage(msg)
		if err == nil {
			return
		}
		s.Logger.Printf("Failed to send message retry count: %d waiting: %s", i, nextWaitDuration)
		time.Sleep(nextWaitDuration)
	}
	err := s.DLQueueWriter.WriteMessage(ctx, msg)
	if err != nil {
		s.Logger.Fatal(err)
	}
}

func (s *Service) ConsumeMessages(ctx context.Context, wg *sync.WaitGroup) {
	messageChan := make(chan queue.MessagePayload, 1000)
	go func() {
		for {
			msg, err := s.RetryQueueReader.ReadMessage(ctx)
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
	}()
	for {
		select {
		case msg, ok := <-messageChan:
			if !ok {
				s.Logger.Println("Message channel closed, stopping consumer")
				wg.Done()
				return
			}
			go s.repeatSendMessage(ctx, msg)
		case <-ctx.Done():
			for {
				select {
				case msg, ok := <-messageChan:
					if !ok {
						wg.Done()
						return
					}
					go s.repeatSendMessage(ctx, msg)
				default:
					wg.Done()
					return
				}
			}
		}
	}

}

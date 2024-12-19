package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"

	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/common/logger"
	"github.com/ozzy-cox/automatic-message-system/internal/common/queue"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	Config            *ConsumerConfig
	Cache             *redis.Client
	MessageRepository db.IMessageRepository
	ReaderQClient     *queue.ReaderClient
	WriterQClient     *queue.WriterClient
	RetryQueueWriter  *queue.WriterClient
	Logger            *logger.Logger
}

func (s *Service) MustSetMessageIdToCache(ctx context.Context, msgId string, msg MessageResponse) error {
	_, err := s.Cache.Set(ctx, msgId, msg, redis.KeepTTL).Result()
	if err != nil {
		s.Logger.Fatalf("Failed to set producer offset in redis: %v", err)
	}
	return nil
}

func (s *Service) RequeueMessage(ctx context.Context, msg queue.MessagePayload) error {
	err := s.RetryQueueWriter.WriteMessage(ctx, msg)
	if err != nil {
		s.Logger.Printf("Error queueing message for retry: %v", err)
	}
	return nil
}

func (s *Service) sendMessage(ctx context.Context, msg queue.MessagePayload) error {
	jsonBody, err := json.Marshal(msg)
	if err != nil {
		s.Logger.Printf("Error marshaling message: %v", err)
		return err
	}
	bodyReader := bytes.NewReader(jsonBody)

	resp, err := http.Post(s.Config.RequestURL, "application/json", bodyReader)
	if err != nil || resp.StatusCode != 200 {
		errMessage := fmt.Sprintf("Error sending message to %s: %v", s.Config.RequestURL, err)
		s.Logger.Println(errMessage)
		return errors.New(errMessage)
	}

	var body MessageResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		s.Logger.Printf("Failed to read response from message response")
	}
	s.Logger.Printf("Successfully sent message ID: %d to %s", msg.ID, msg.To)

	err = s.MessageRepository.SetMessageSent(msg.ID)
	if err != nil {
		s.Logger.Printf("Failed to update message sent state: %v", err)
		return err
	}

	if body.MessageId != nil {
		s.MustSetMessageIdToCache(ctx, *body.MessageId, body)
	}
	return nil
}

func (s *Service) handleMessage(ctx context.Context, msg queue.MessagePayload) {
	for i := 1; i < s.Config.RetryCount; i++ {
		jitter := (rand.Float64() * 1) - 1
		nextWaitDuration := time.Duration(math.Pow(2, float64(i))+jitter) * s.Config.Interval
		err := s.sendMessage(ctx, msg)
		if err == nil {
			return
		}
		s.Logger.Printf("Failed to send message retry count: %d waiting: %s", i, nextWaitDuration)
		time.Sleep(nextWaitDuration)
	}
	go s.RequeueMessage(ctx, msg)
}

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
	for {
		select {
		case msg, ok := <-messageChan:
			if !ok {
				s.Logger.Println("Message channel closed, stopping consumer")
				wg.Done()
				return
			}
			go s.handleMessage(ctx, msg)
		case <-ctx.Done():
			for {
				select {
				case msg, ok := <-messageChan:
					if !ok {
						wg.Done()
						return
					}
					go s.handleMessage(ctx, msg)
				default:
					wg.Done()
					return
				}
			}
		}
	}

}

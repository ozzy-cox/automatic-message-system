package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/common/logger"
	"github.com/ozzy-cox/automatic-message-system/internal/common/queue"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	Config            *ConsumerConfig
	Cache             *redis.Client
	MessageRepository db.IMessageRepository
	QueueReader       *queue.ReaderClient
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

func (s *Service) TryRequeueRetryMessage(ctx context.Context, msg queue.MessagePayload) error {
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
		s.Logger.Printf("Error sending message to %s: %v", s.Config.RequestURL, err)
		return err
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

func (s *Service) ConsumeMessages(ctx context.Context, wg *sync.WaitGroup) {
	messageChan := make(chan queue.MessagePayload, 1000)
	go func() {
		for {
			msg, err := s.QueueReader.ReadMessage(ctx)
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
			err := s.sendMessage(ctx, msg)
			if err != nil {
				go s.TryRequeueRetryMessage(ctx, msg)
			}
		case <-ctx.Done():
			for {
				select {
				case msg, ok := <-messageChan:
					if !ok {
						wg.Done()
						return
					}
					err := s.sendMessage(ctx, msg)
					if err != nil {
						s.TryRequeueRetryMessage(ctx, msg)
					}
				default:
					wg.Done()
					return
				}
			}
		}
	}

}

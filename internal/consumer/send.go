package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/ozzy-cox/automatic-message-system/internal/common/queue"
)

func (s *Service) handleMessage(ctx context.Context, msg queue.MessagePayload) {
	for i := 1; i < s.Config.RetryCount; i++ {
		jitter := (rand.Float64() * 1) - 1
		nextWaitDuration := time.Duration(math.Pow(2, float64(i))+jitter) * s.Config.Interval
		err := s.sendMessage(ctx, msg)
		if err == nil {
			s.Logger.Printf("Successfully sent message ID: %d to %s", msg.ID, msg.To)
			return
		}
		s.Logger.Printf("Failed to send message retry count: %d waiting: %s cause: %v", i, nextWaitDuration, err)
		time.Sleep(nextWaitDuration)
	}
	go s.requeueMessage(ctx, msg)
}

func (s *Service) sendMessage(ctx context.Context, msg queue.MessagePayload) error {
	jsonBody, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(jsonBody)

	resp, err := http.Post(s.Config.RequestURL, "application/json", bodyReader)
	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("Error sending message to %s: %v", s.Config.RequestURL, err)
	}

	var body MessageResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return err
	}

	err = s.MessageRepository.MarkMessageAsSent(msg.ID)
	if err != nil {
		return err
	}

	if body.MessageId != nil {
		s.cacheMessageResponse(ctx, *body.MessageId, body)
	}
	return nil
}

func (s *Service) requeueMessage(ctx context.Context, msg queue.MessagePayload) error {
	err := s.RetryQueueWriter.WriteMessages(ctx, msg)
	if err != nil {
		s.Logger.Printf("Error queueing message for retry: %v", err)
		return err
	}
	return nil
}

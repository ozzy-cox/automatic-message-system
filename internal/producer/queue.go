package producer

import (
	"context"
	"iter"

	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/common/queue"
)

func (s *Service) parseMessages(messages iter.Seq2[*db.Message, error]) []queue.MessagePayload {
	parsedMessages := make([]queue.MessagePayload, 0)
	for msg, err := range messages {
		if err != nil {
			s.Logger.Printf("Error scanning messages: %v", err)
			continue
		}
		parsedMessages = append(parsedMessages, queue.MessagePayload{
			ID:        msg.ID,
			Content:   msg.Content,
			To:        msg.To,
			CreatedAt: msg.CreatedAt,
		})
		s.Logger.Printf("Successfully queued message ID: %d for recipient: %s", msg.ID, msg.To)
	}
	return parsedMessages
}

func (s *Service) PushMessagesToQ(ctx context.Context, limit, offset int) int {
	s.Logger.Printf("Fetching messages starting at offset: %d", offset)
	messages := s.MessageRepository.GetMessages(limit, offset)
	parsedMessages := s.parseMessages(messages)

	if err := s.Queue.WriteMessages(ctx, parsedMessages...); err != nil {
		s.Logger.Printf("Error writing messages to queue: %v", err)
	}
	return len(parsedMessages)
}

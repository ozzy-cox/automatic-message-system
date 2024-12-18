package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

const (
	TopicMessages = "messages"
)

type WriterClient struct {
	writer *kafka.Writer
}
type ReaderClient struct {
	reader *kafka.Reader
}

func NewReaderClient(cfg KafkaConfig) (*ReaderClient, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		Topic:    TopicMessages,
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	return &ReaderClient{reader: reader}, nil
}

func NewWriterClient(cfg KafkaConfig) (*WriterClient, error) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Brokers,
		Topic:   TopicMessages,
	})

	return &WriterClient{writer: writer}, nil
}

func (c *WriterClient) Close() error {
	if err := c.writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}
	return nil
}

func (c *ReaderClient) Close() error {
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("failed to close reader: %w", err)
	}
	return nil
}

func (c *WriterClient) WriteMessage(ctx context.Context, msg MessagePayload) error {
	value, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// TODO dont serialize messages here
	return c.writer.WriteMessages(ctx, kafka.Message{
		Value: value,
	})
}

func (c *ReaderClient) ReadMessage(ctx context.Context) (MessagePayload, error) {
	var message MessagePayload

	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return message, fmt.Errorf("failed to read message: %w", err)
	}

	if err := json.Unmarshal(msg.Value, &message); err != nil {
		return message, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return message, nil
}
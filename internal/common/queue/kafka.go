package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaWriterClient struct {
	writer *kafka.Writer
}

type KafkaReaderClient struct {
	reader *kafka.Reader
}

func NewKafkaReaderClient(cfg KafkaConfig) (*KafkaReaderClient, error) {
	conn := mustEnsureConn(cfg)
	tryCreateTopic(conn, cfg.Topic)
	defer conn.Close()
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		GroupID:  cfg.GroupID,
		Topic:    cfg.Topic,
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	return &KafkaReaderClient{reader: reader}, nil
}

func NewKafkaWriterClient(cfg KafkaConfig) (*KafkaWriterClient, error) {
	conn := mustEnsureConn(cfg)
	tryCreateTopic(conn, cfg.Topic)
	defer conn.Close()
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Brokers...),
		Topic:                  cfg.Topic,
		AllowAutoTopicCreation: true,
		BatchTimeout:           time.Millisecond,
	}

	return &KafkaWriterClient{writer: writer}, nil
}

func (c *KafkaWriterClient) Close() error {
	if err := c.writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}
	return nil
}

func (c *KafkaReaderClient) Close() error {
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("failed to close reader: %w", err)
	}
	return nil
}

func (c *KafkaWriterClient) WriteMessages(ctx context.Context, msgs ...MessagePayload) error {
	kafkaMessages := make([]kafka.Message, len(msgs))
	for i, msg := range msgs {
		value, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}
		kafkaMessages[i] = kafka.Message{
			Value: value,
		}
	}

	return c.writer.WriteMessages(ctx, kafkaMessages...)
}

func (c *KafkaReaderClient) ReadMessage(ctx context.Context) (MessagePayload, error) {
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

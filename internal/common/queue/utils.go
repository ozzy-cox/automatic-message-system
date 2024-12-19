package queue

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

func mustEnsureConn(cfg KafkaConfig) *kafka.Conn {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := kafka.DialLeader(ctx, "tcp", cfg.Brokers[0], cfg.Topic, 0)
	if err != nil {
		log.Fatalf("Failed to connect to leader: %v", err)
	}
	return conn
}

func tryCreateTopic(conn *kafka.Conn, topic string) {
	controller, err := conn.Controller()
	if err != nil {
		log.Fatalf("failed to get controller: %v", err)
	}
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Fatalf("failed to connect to controller: %v", err)
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	if err != nil {
		log.Printf("Error creating topic: %v", err)
	}
}

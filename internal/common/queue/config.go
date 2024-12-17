package queue

import (
	"strings"

	"github.com/ozzy-cox/automatic-message-system/internal/common/utils"
)

type KafkaConfig struct {
	Brokers []string
	GroupID string
	Topic   string
}

func GetDefaultKafkaBrokers() []string {
	brokersStr := utils.GetEnvStringWithDefault("KAFKA_BROKERS", "localhost:9092")
	brokers := strings.Split(brokersStr, ",")
	return brokers
}

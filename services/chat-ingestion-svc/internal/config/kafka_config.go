package config

import (
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/helper/logs"
	"be-realtime-chat-app/services/commoner/utils"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

func NewKafkaProducer(log logs.Log) *kafka.Producer {
	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers": utils.GetEnv("KAFKA_BOOTSTRAP_SERVERS"),
	}

	producer, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		log.Fatal("Failed to create producer: %v", zap.Error(err))
	}
	return producer
}

package config

import (
	"be-realtime-chat-app/services/commoner/logs"
	"be-realtime-chat-app/services/commoner/utils"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

func NewKafkaConsumer(log logs.Log) *kafka.Consumer {
	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers":     utils.GetEnv("KAFKA_BOOTSTRAP_SERVERS"),
		"group.id":              utils.GetEnv("KAFKA_GROUP_ID"),
		"auto.offset.reset":     utils.GetEnv("KAFKA_AUTO_OFFSET_RESET"),
		"broker.address.family": utils.GetEnv("BROKER_ADDRESS_FAMILY"),
	}

	consumer, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		log.Fatal("Failed to create consumer:", zap.Error(err))
	}
	return consumer
}

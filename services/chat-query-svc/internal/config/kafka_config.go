package config

import (
	"be-realtime-chat-app/services/commoner/utils"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewKafkaProducer(config *viper.Viper, log *logrus.Logger) *kafka.Producer {
	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers": utils.GetEnv("KAFKA_BOOTSTRAP_SERVERS"),
	}

	producer, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	return producer
}

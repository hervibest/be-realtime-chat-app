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

		// No need for EOS
		"enable.idempotence": false,

		// Fan-out friendly: faster & allows batching
		"acks":                                  "1", // Cepat, hanya perlu 1 broker ack
		"retries":                               3,
		"linger.ms":                             5,         // Kumpulkan pesan dalam 5ms
		"batch.size":                            64 * 1024, // 64KB batching
		"compression.type":                      "snappy",  // Ringan dan cepat
		"max.in.flight.requests.per.connection": 10,

		// Optional buffer tuning
		"queue.buffering.max.messages": 100000,
		"queue.buffering.max.kbytes":   1048576, // 1GB
	}

	producer, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		log.Fatal("Failed to create producer: %v", zap.Error(err))
	}
	return producer
}

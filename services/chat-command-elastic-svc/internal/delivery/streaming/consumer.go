package consumer

import (
	"be-realtime-chat-app/services/commoner/logs"
	"context"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

type ConsumerHandler func(message *kafka.Message) error

func ConsumeTopic(ctx context.Context, consumer *kafka.Consumer, topic string, log logs.Log, handler ConsumerHandler) {
	err := consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatal("Failed to subscribe to topic: %v", zap.Error(err), zap.String("topic", topic))
		return
	}

	run := true

	for run {
		select {
		case <-ctx.Done():
			run = false
		default:
			message, err := consumer.ReadMessage(time.Second)
			if err == nil {
				err := handler(message)
				if err != nil {
					log.Error("Error processing message: %v", zap.Error(err), zap.String("topic", topic), zap.Int32("partition", message.TopicPartition.Partition))
				} else {
					_, err = consumer.CommitMessage(message)
					if err != nil {
						log.Error("Failed to commit message: %v", zap.Error(err))
					} else {
						log.Info("Message committed successfully", zap.String("topic", topic), zap.Int32("partition", message.TopicPartition.Partition))
					}
				}
			} else if !err.(kafka.Error).IsTimeout() {
				log.Error("Error reading message: %v", zap.Error(err))
				run = false
			} else {
				log.Debug("No message received, waiting for next message")
			}
		}
	}

	log.Info("Closing consumer for topic :", zap.String("topic", topic))
	err = consumer.Close()
	if err != nil {
		panic(err)
	}
}

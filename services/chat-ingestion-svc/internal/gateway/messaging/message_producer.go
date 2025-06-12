package producer

import (
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/adapter"
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/model/event"
	"be-realtime-chat-app/services/commoner/logs"

	"github.com/bytedance/sonic"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"go.uber.org/zap"
)

type MessageProducer interface {
	ProduceMessage(event *event.Message) error
}
type messageProducerImpl struct {
	streaming adapter.StreamingAdapter
	log       logs.Log
}

func NewMessageProducer(streaming adapter.StreamingAdapter, log logs.Log) MessageProducer {
	return &messageProducerImpl{streaming: streaming, log: log}
}

func (p *messageProducerImpl) ProduceMessage(event *event.Message) error {
	p.log.Info("Processing message", zap.String("roomID", event.RoomID), zap.String("message", string(event.Content)))
	value, err := sonic.ConfigFastest.Marshal(event)
	if err != nil {
		p.log.Error("Failed to marshal message", zap.Error(err), zap.String("message", event.RoomID))
		return err
	}

	topic := "room"
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: value,
		Key:   []byte(event.GetID()),
	}

	if err = p.streaming.Produce(message, nil); err != nil {
		p.log.Error("Failed to produce message to Kafka", zap.Error(err), zap.String("topic", topic), zap.String("message", string(value)))
		return err
	}
	return nil
}

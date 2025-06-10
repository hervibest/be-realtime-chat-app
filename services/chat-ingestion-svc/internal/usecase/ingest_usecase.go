package usecase

import (
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/adapter"
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/helper/logs"
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/model/event"
	"context"
	"log"

	"github.com/bytedance/sonic"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type RoomUseCase interface {
	IngestChat(ctx context.Context) error
}

type roomUseCaseImpl struct {
	mesaging  adapter.Messaging
	streaming adapter.Streaming
	log       logs.Log
}

func NewRoomUseCase(mesaging adapter.Messaging, streaming adapter.Streaming, log logs.Log) RoomUseCase {
	return &roomUseCaseImpl{
		mesaging:  mesaging,
		streaming: streaming,
		log:       log,
	}
}

func (uc *roomUseCaseImpl) IngestChat(ctx context.Context) error {
	bufferedMsg := make(chan *event.Message, 20)
	go func() {
		for event := range bufferedMsg {
			if event == nil {
				uc.log.Warn("Received nil message in buffered channel")
				continue
			}

			uc.log.Info("Processing message", zap.String("roomID", event.RoomID), zap.String("message", string(event.Content)))
			value, err := sonic.ConfigFastest.Marshal(event)
			if err != nil {
				uc.log.Error("Failed to marshal message", zap.Error(err), zap.String("message", event.RoomID))
				continue
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

			if err = uc.streaming.Produce(message, nil); err != nil {
				uc.log.Error("Failed to produce message to Kafka", zap.Error(err), zap.String("topic", topic), zap.String("message", string(value)))
				continue
			}
		}
	}()

	sub, err := uc.mesaging.Subscribe("room.", func(msg *nats.Msg) {
		if msg == nil {
			uc.log.Warn("Received nil message")
			return
		}

		event := new(event.Message)
		if err := sonic.ConfigFastest.Unmarshal(msg.Data, event); err != nil {
			uc.log.Error("Failed to unmarshal message", zap.Error(err), zap.String("data", string(msg.Data)))
			return
		}

		uc.log.Info("Received message for room", zap.String("roomID", event.RoomID), zap.String("message", string(event.Content)))
		bufferedMsg <- event
	})

	if err != nil {
		log.Println("NATS Subscribe error:", err)
		return err
	}

	sub.Unsubscribe()
	defer close(bufferedMsg)
	return nil
}

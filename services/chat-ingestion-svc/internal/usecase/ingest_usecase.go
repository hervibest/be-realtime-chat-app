package usecase

import (
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/adapter"
	producer "be-realtime-chat-app/services/chat-ingestion-svc/internal/gateway/messaging"
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/helper/logs"
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/model/event"
	"context"

	"go.uber.org/zap"
)

type IngestUeCase interface {
	IngestChat(ctx context.Context, bufferedMsg chan *event.Message)
}

type roomUseCaseImpl struct {
	mesaging adapter.Messaging
	producer producer.MessageProducer
	log      logs.Log
}

func NewIngestUseCase(mesaging adapter.Messaging, producer producer.MessageProducer, log logs.Log) IngestUeCase {
	return &roomUseCaseImpl{
		mesaging: mesaging,
		producer: producer,
		log:      log,
	}
}

func (uc *roomUseCaseImpl) IngestChat(ctx context.Context, bufferedMsg chan *event.Message) {
	uc.log.Info("Starting chat ingestion process")
	for event := range bufferedMsg {
		if event == nil {
			uc.log.Warn("Received nil message in buffered channel")
			continue
		}

		if err := uc.producer.ProduceMessage(event); err != nil {
			uc.log.Error("Failed to produce message", zap.Error(err), zap.String("roomID", event.RoomID), zap.String("message", string(event.Content)))
			continue
		}

	}

}

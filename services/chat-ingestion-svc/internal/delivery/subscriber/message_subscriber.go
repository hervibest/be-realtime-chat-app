package subscriber

import (
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/adapter"
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/model/event"
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/usecase"
	"be-realtime-chat-app/services/commoner/logs"
	"context"
	"log"

	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type MessageSubscriber interface {
	SubscribeToMessages() error
}
type messageSubscriberImpl struct {
	ingestUseCase usecase.IngestUeCase
	messaging     adapter.Messaging
	log           logs.Log
}

func NewMessageSubscriber(ingestUseCase usecase.IngestUeCase, messaging adapter.Messaging, log logs.Log) MessageSubscriber {
	return &messageSubscriberImpl{ingestUseCase: ingestUseCase, messaging: messaging, log: log}
}

// TODO GRACEFUL TOPIC ENVARS INJECTION
func (s *messageSubscriberImpl) SubscribeToMessages() error {
	bufferedMsg := make(chan *event.Message, 100)
	s.log.Info("Starting message subscription")
	go s.ingestUseCase.IngestChat(context.Background(), bufferedMsg)
	sub, err := s.messaging.Subscribe("room.", func(msg *nats.Msg) {
		if msg == nil {
			s.log.Warn("Received nil message")
			return
		}

		event := new(event.Message)
		if err := sonic.ConfigFastest.Unmarshal(msg.Data, event); err != nil {
			s.log.Error("Failed to unmarshal message", zap.Error(err), zap.String("data", string(msg.Data)))
			return
		}

		s.log.Info("Received message for room", zap.String("roomID", event.RoomID), zap.String("message", string(event.Content)))
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

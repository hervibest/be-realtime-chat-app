package subscriber

import (
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/adapter"
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/model/event"
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/usecase"
	"be-realtime-chat-app/services/commoner/logs"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type MessageSubscriber interface {
	SubscribeToMessages(ctx context.Context) error
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
func (s *messageSubscriberImpl) SubscribeToMessages(ctx context.Context) error {
	bufferedMsg := make(chan *event.Message, 100)
	s.log.Info("Starting message subscription")
	go s.ingestUseCase.IngestChat(ctx, bufferedMsg)
	sub, err := s.messaging.SubscribeSync("room.*")
	if err != nil {
		log.Println("NATS Subscribe error:", err)
		return err
	}

	defer func() {
		close(bufferedMsg)
		sub.Unsubscribe()
	}()

	for {
		select {
		case <-ctx.Done():
			s.log.Info("Context done, exiting message subscription")
			return nil
		default:
			// Tunggu pesan NATS dengan timeout
			msg, err := sub.NextMsg(10 * time.Second)
			if err != nil {
				if err == nats.ErrTimeout {
					continue // tidak ada pesan, ulangi
				}
				log.Println("Error receiving message from NATS:", err)
				return nil
			}

			event := new(event.Message)
			if err := sonic.ConfigFastest.Unmarshal(msg.Data, event); err != nil {
				s.log.Error("Failed to unmarshal message", zap.Error(err), zap.String("data", string(msg.Data)))
				return nil
			}

			s.log.Info("Subs called", zap.String("message_id", event.ID), zap.String("created_at", event.CreatedAt))
			createdAt, err := time.Parse(time.RFC3339Nano, event.CreatedAt)
			if err != nil {
				s.log.Error("error parsing created_at", zap.Error(err), zap.String("created_at", event.CreatedAt))
				return fmt.Errorf("invalid created_at format: %w", err)
			}

			duration := time.Since(createdAt)
			duration.Seconds()
			s.log.Info("Subs Received  at", zap.String("created_at", event.CreatedAt), zap.Float64("duration_since_created", duration.Seconds()))

			s.log.Info("Received message from NATS", zap.String("roomID", event.RoomID), zap.String("message", string(event.Content)))

			s.log.Info("Received message for room created at", zap.String("roomID", event.RoomID), zap.String("created_at", event.CreatedAt))
			bufferedMsg <- event
		}
	}

}

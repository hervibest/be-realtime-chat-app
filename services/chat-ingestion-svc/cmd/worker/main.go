package main

import (
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/config"
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/delivery/subscriber"
	producer "be-realtime-chat-app/services/chat-ingestion-svc/internal/gateway/messaging"
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/usecase"
	"be-realtime-chat-app/services/commoner/logs"
	"context"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func worker(ctx context.Context) error {
	logger, _ := logs.NewLogger()
	logger.Info("Starting chat ingestion service...")
	kafkaProducer := config.NewKafkaProducer(logger)
	natsConn := config.NewNATsConn(logger)

	messageProducer := producer.NewMessageProducer(kafkaProducer, logger)
	ingestUseCase := usecase.NewIngestUseCase(natsConn, messageProducer, logger)
	messageSubscriber := subscriber.NewMessageSubscriber(ingestUseCase, natsConn, logger)
	if err := messageSubscriber.SubscribeToMessages(ctx); err != nil {
		logger.Fatal("Failed to subscribe to message topic", zap.Error(err))
	}
	return nil

}
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := worker(ctx); err != nil {
	}
}

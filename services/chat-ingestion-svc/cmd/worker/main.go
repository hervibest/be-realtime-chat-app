package main

import (
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/config"
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/delivery/subscriber"
	producer "be-realtime-chat-app/services/chat-ingestion-svc/internal/gateway/messaging"
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/usecase"
	"be-realtime-chat-app/services/commoner/logs"

	"go.uber.org/zap"
)

func main() {
	logger, _ := logs.NewLogger()
	kafkaProducer := config.NewKafkaProducer(logger)
	natsConn := config.NewNATsConn(logger)

	messageProducer := producer.NewMessageProducer(kafkaProducer, logger)
	ingestUseCase := usecase.NewIngestUseCase(natsConn, messageProducer, logger)
	messageSubscriber := subscriber.NewMessageSubscriber(ingestUseCase, natsConn, logger)
	if err := messageSubscriber.SubscribeToMessages(); err != nil {
		logger.Fatal("Failed to subscribe to message topic", zap.Error(err))
	}
}

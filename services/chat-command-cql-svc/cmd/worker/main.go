package main

import (
	"be-realtime-chat-app/services/chat-command-cql-svc/internal/config"
	consumer "be-realtime-chat-app/services/chat-command-cql-svc/internal/delivery/streaming"
	"be-realtime-chat-app/services/chat-command-cql-svc/internal/repository"
	"be-realtime-chat-app/services/chat-command-cql-svc/internal/usecase"
	"be-realtime-chat-app/services/commoner/logs"
	"context"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger, err := logs.NewLogger()
	if err != nil {
		panic(err)
	}

	kafkaConsumer := config.NewKafkaConsumer(logger)
	cqlClient, err := config.NewCQLDB()

	messageRepository := repository.NewMessageRepository(cqlClient)
	commandAsyncUseCase := usecase.NewCommandUseCase(messageRepository, logger)
	messageConsumer := consumer.NewMessageConsumerImpl(commandAsyncUseCase, logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go consumer.ConsumeTopic(ctx, kafkaConsumer, "room", logger, messageConsumer.Consume)

	logger.Info("Chat Command CQL Service Worker started")
	<-ctx.Done()

	time.Sleep(2 * time.Second) // Allow some time for graceful shutdown
	logger.Info("Chat Command CQL Service Worker stopped")
}

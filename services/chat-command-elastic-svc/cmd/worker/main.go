package main

import (
	"be-realtime-chat-app/services/chat-command-elastic-svc/internal/config"
	consumer "be-realtime-chat-app/services/chat-command-elastic-svc/internal/delivery/streaming"
	"be-realtime-chat-app/services/chat-command-elastic-svc/internal/repository"
	"be-realtime-chat-app/services/chat-command-elastic-svc/internal/usecase"
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
	elasticsearchClient, err := config.NewElasticsearch()
	// if err := config.CreateMessageIndex(elasticsearchClient); err != nil {
	// 	logger.Fatal("Failed to create message index", zap.Error(err))
	// }

	messageRepository := repository.NewMessageRepository(elasticsearchClient)
	commandAsyncUseCase := usecase.NewCommandAsyncUseCase(messageRepository, logger)
	messageConsumer := consumer.NewMessageConsumerImpl(commandAsyncUseCase, logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go consumer.ConsumeTopic(ctx, kafkaConsumer, "room", logger, messageConsumer.Consume)

	logger.Info("Chat Command Elastic Service Worker started")
	<-ctx.Done()

	time.Sleep(2 * time.Second) // Allow some time for graceful shutdown
	logger.Info("Chat Command Elastic Service Worker stopped")
}

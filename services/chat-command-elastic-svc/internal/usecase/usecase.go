package usecase

import (
	"be-realtime-chat-app/services/chat-command-elastic-svc/internal/adapter"
	"be-realtime-chat-app/services/chat-command-elastic-svc/internal/helper/logs"
	"be-realtime-chat-app/services/chat-command-elastic-svc/internal/model/event"
	"be-realtime-chat-app/services/chat-command-elastic-svc/internal/repository"
	"context"
)

type CommandAsyncUseCase interface {
	Persist(ctx context.Context, request *event.Message) error
}

type messageUseCaseImpl struct {
	messageRepository repository.MessageRepository
	db                repository.DB
	log               logs.Log
}

func NewCommandAsyncUseCase(messageRepository repository.MessageRepository, db repository.DB, mesaging adapter.Messaging,
	streaming adapter.Streaming, log logs.Log) CommandAsyncUseCase {
	return &messageUseCaseImpl{
		messageRepository: messageRepository,
		db:                db,
		log:               log,
	}
}

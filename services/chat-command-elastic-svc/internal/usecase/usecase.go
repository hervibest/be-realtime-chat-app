package usecase

import (
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
	log               logs.Log
}

func NewCommandAsyncUseCase(messageRepository repository.MessageRepository, log logs.Log) CommandAsyncUseCase {
	return &messageUseCaseImpl{
		messageRepository: messageRepository,
		log:               log,
	}
}

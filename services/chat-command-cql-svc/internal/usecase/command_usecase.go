package usecase

import (
	"be-realtime-chat-app/services/chat-command-cql-svc/internal/entity"
	"be-realtime-chat-app/services/chat-command-cql-svc/internal/helper/logs"
	"be-realtime-chat-app/services/chat-command-cql-svc/internal/model/event"
	"be-realtime-chat-app/services/chat-command-cql-svc/internal/repository"
	"context"
	"time"
)

type CommandUseCase interface {
	PersistChat(ctx context.Context, request *event.Message) error
}

type commandUseCaseImpl struct {
	messageRepository repository.MessageRepository
	log               logs.Log
}

func NewCommandUseCase(messageRepository repository.MessageRepository, log logs.Log) CommandUseCase {
	return &commandUseCaseImpl{
		messageRepository: messageRepository,
		log:               log,
	}
}

func (uc *commandUseCaseImpl) PersistChat(ctx context.Context, request *event.Message) error {
	createdAt, _ := time.Parse(time.RFC3339, request.CreatedAt)
	message := &entity.Message{
		ID:        request.ID,
		RoomID:    request.RoomID,
		UserID:    request.UserID,
		Content:   request.Content,
		CreatedAt: createdAt,
	}

	if err := uc.messageRepository.Insert(message); err != nil {
		return err
	}

	return nil
}

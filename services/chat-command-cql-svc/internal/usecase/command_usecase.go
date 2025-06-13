package usecase

import (
	"be-realtime-chat-app/services/chat-command-cql-svc/internal/entity"
	"be-realtime-chat-app/services/chat-command-cql-svc/internal/helper/logs"
	"be-realtime-chat-app/services/chat-command-cql-svc/internal/model/event"
	"be-realtime-chat-app/services/chat-command-cql-svc/internal/repository"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
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
	uc.log.Info("PersistChat called", zap.String("message_id", request.ID), zap.String("created_at", request.CreatedAt))
	createdAt, err := time.Parse(time.RFC3339Nano, request.CreatedAt)
	if err != nil {
		uc.log.Error("error parsing created_at", zap.Error(err), zap.String("created_at", request.CreatedAt))
		return fmt.Errorf("invalid created_at format: %w", err)
	}

	duration := time.Since(createdAt)
	duration.Seconds()
	uc.log.Info("Message created at", zap.String("created_at", request.CreatedAt), zap.Float64("duration_since_created", duration.Seconds()))

	message := &entity.Message{
		ID:        request.ID,
		UUID:      request.UUID,
		RoomID:    request.RoomID,
		UserID:    request.UserID,
		Username:  request.Username,
		Content:   request.Content,
		CreatedAt: createdAt,
	}

	fmt.Println("Persisting message:", message)

	if err := uc.messageRepository.Insert(message); err != nil {
		return err
	}

	return nil
}

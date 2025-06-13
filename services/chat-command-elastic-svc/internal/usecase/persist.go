package usecase

import (
	"be-realtime-chat-app/services/chat-command-elastic-svc/internal/entity"
	"be-realtime-chat-app/services/chat-command-elastic-svc/internal/model/event"
	"context"
	"time"
)

func (uc *messageUseCaseImpl) Persist(ctx context.Context, request *event.Message) error {
	createdAt, _ := time.Parse(time.RFC3339Nano, request.CreatedAt)
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

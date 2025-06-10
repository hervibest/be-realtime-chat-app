package converter

import (
	"be-realtime-chat-app/services/chat-realtime-svc/internal/entity"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/model"
	"time"
)

func RoomToResponse(room *entity.Room) *model.RoomResponse {
	if room == nil {
		return nil
	}

	return &model.RoomResponse{
		ID:        room.ID,
		UUID:      room.UUID,
		Name:      room.Name,
		UserID:    room.UserID,
		CreatedAt: room.CreatedAt.Format(time.RFC3339),
	}
}

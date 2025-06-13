package converter

import (
	"be-realtime-chat-app/services/room-svc/internal/entity"
	"be-realtime-chat-app/services/room-svc/internal/model"
)

func RoomToResponse(room *entity.Room) *model.RoomResponse {
	if room == nil {
		return nil
	}

	return &model.RoomResponse{
		ID:     room.ID,
		UUID:   room.UUID,
		Name:   room.Name,
		UserID: room.UserID,
		// CreatedAt: room.CreatedAt.Format(time.RFC3339Nano,
	}
}

func RoomsToResponses(rooms *[]*entity.Room) *[]*model.RoomResponse {
	if rooms == nil {
		return nil
	}

	responses := make([]*model.RoomResponse, len(*rooms))
	for i, room := range *rooms {
		responses[i] = RoomToResponse(room)
	}

	return &responses
}

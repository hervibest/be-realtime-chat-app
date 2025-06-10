package converter

import (
	"be-realtime-chat-app/services/chat-query-svc/internal/entity"
	"be-realtime-chat-app/services/chat-query-svc/internal/model"
)

func MessageToResponse(message *entity.Message) *model.MessageResponse {
	if message == nil {
		return nil
	}
	return &model.MessageResponse{
		ID:        message.ID,
		RoomID:    message.RoomID,
		UserID:    message.UserID,
		Content:   message.Content,
		CreatedAt: message.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func MessagesToResponses(messages *[]*entity.Message) *[]*model.MessageResponse {
	if *messages == nil {
		return nil
	}

	responses := make([]*model.MessageResponse, len(*messages))
	for i, message := range *messages {
		responses[i] = MessageToResponse(message)
	}
	return &responses
}

package converter

import (
	"be-realtime-chat-app/services/user-svc/internal/entity"
	"be-realtime-chat-app/services/user-svc/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	if user == nil {
		return nil
	}

	return &model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}

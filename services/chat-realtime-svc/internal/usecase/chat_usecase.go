package usecase

import (
	"be-realtime-chat-app/services/chat-realtime-svc/internal/adapter"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/helper/logs"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/model"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/websockets"
	errorcode "be-realtime-chat-app/services/commoner/constant/errcode"
	"be-realtime-chat-app/services/commoner/helper"
	"context"

	"github.com/gofiber/contrib/websocket"
	"go.uber.org/zap"
)

type ChatUseCase interface {
	JoinRoom(ctx context.Context, conn *websocket.Conn, request *model.JoinRoomRequest) error
}

type chatUseCase struct {
	messagingAdapter adapter.MessagingAdapter
	roomAdapter      adapter.RoomAdapter
	customValidator  helper.CustomValidator
	log              logs.Log
}

func NewChatUseCase(messagingAdapter adapter.MessagingAdapter, roomAdapter adapter.RoomAdapter,
	customValidator helper.CustomValidator, log logs.Log) ChatUseCase {
	return &chatUseCase{
		messagingAdapter: messagingAdapter,
		roomAdapter:      roomAdapter,
		customValidator:  customValidator,
		log:              log,
	}
}

// TODO : add cache for room repo
func (uc *chatUseCase) JoinRoom(ctx context.Context, conn *websocket.Conn, request *model.JoinRoomRequest) error {
	if validatonErrs := uc.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return validatonErrs
	}

	uc.log.Info("Join room request", zap.String("roomID", request.RoomID), zap.String("roomID", request.UserID))
	room, err := uc.roomAdapter.GetRoom(ctx, request.RoomID)
	if err != nil {
		if appErr, ok := err.(*helper.AppError); ok {
			if appErr.Code == errorcode.ErrUserNotFound {
				return helper.NewUseCaseError(errorcode.ErrUserNotFound, "Room not found")
			}
			return appErr
		}
	}

	client := &websockets.UserClient{
		Conn:      conn,
		Messaging: uc.messagingAdapter,
		UserID:    request.UserID,
		RoomID:    room.Room.GetId(),
		Email:     request.Email,
	}

	go client.Subscriber()
	client.Publisher()

	return nil
}

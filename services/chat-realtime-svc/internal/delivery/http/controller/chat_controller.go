package controller

import (
	"be-realtime-chat-app/services/chat-realtime-svc/internal/delivery/http/middleware"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/helper/logs"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/model"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/usecase"
	"be-realtime-chat-app/services/commoner/helper"
	"net/http"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

type ChatController interface {
	JoinRoom(ctx *fiber.Ctx) error
}

type roomControllerImpl struct {
	roomUseCase usecase.ChatUseCase
	log         logs.Log
}

func NewChatController(roomUseCase usecase.ChatUseCase, log logs.Log) ChatController {
	return &roomControllerImpl{
		roomUseCase: roomUseCase,
		log:         log,
	}
}

func (c *roomControllerImpl) JoinRoom(ctx *fiber.Ctx) error {
	request := new(model.JoinRoomRequest)
	roomID := ctx.Params("roomID")
	_, err := ulid.Parse(roomID)
	if err != nil {
		return helper.ErrCustomResponseJSON(ctx, http.StatusBadRequest, "Invalid room ID format")
	}

	request.RoomID = roomID
	user := middleware.GetUser(ctx)
	request.UserID = user.ID
	request.Email = user.Email
	request.Username = user.Username

	if websocket.IsWebSocketUpgrade(ctx) {
		return websocket.New(func(conn *websocket.Conn) {
			if err := c.roomUseCase.JoinRoom(ctx.Context(), conn, request); err != nil {
				return
			}
		})(ctx)
	}

	return helper.ErrCustomResponseJSON(ctx, http.StatusUpgradeRequired, "WebSocket upgrade required")
}

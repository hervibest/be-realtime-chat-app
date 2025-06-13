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
	"go.uber.org/zap"
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
	c.log.Info("JoinRoom called", zap.String("roomID", ctx.Params("roomID")))
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
		// Capture the context before entering the WebSocket handler
		parentCtx := ctx.Context()

		return websocket.New(func(conn *websocket.Conn) {
			// Use the captured context
			if err := c.roomUseCase.JoinRoom(parentCtx, conn, request); err != nil {
				if appErr, ok := err.(*helper.AppError); ok {
					c.log.Info(appErr.Message)
				}
				return
			}
		})(ctx)
	}

	return helper.ErrCustomResponseJSON(ctx, http.StatusUpgradeRequired, "WebSocket upgrade required")
}

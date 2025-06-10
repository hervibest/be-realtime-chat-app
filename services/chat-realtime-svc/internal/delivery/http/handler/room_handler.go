package handler

import (
	"be-realtime-chat-app/services/chat-realtime-svc/internal/model"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/usecase"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type RoomHandler interface {
	CreateRoom(ctx *fiber.Ctx) error
	JoinRoom(ctx *fiber.Ctx) error
}

type roomHandlerImpl struct {
	roomUseCase usecase.RoomUseCase
}

func NewRoomHandler(roomUseCase usecase.RoomUseCase) RoomHandler {
	return &roomHandlerImpl{
		roomUseCase: roomUseCase,
	}
}

func (r *roomHandlerImpl) CreateRoom(ctx *fiber.Ctx) error {
	request := new(model.CreateRoomRequest)
	if err := ctx.BodyParser(request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if request.Name == "" || request.UserID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Room name and User ID are required",
		})
	}

	response, err := r.roomUseCase.CreateRoom(ctx.Context(), request)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create room: " + err.Error(),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (r *roomHandlerImpl) JoinRoom(ctx *fiber.Ctx) error {
	roomID := ctx.Params("room")
	userID := ctx.Params("user")

	if roomID == "" || userID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Room ID and User ID are required",
		})
	}

	if websocket.IsWebSocketUpgrade(ctx) {
		return websocket.New(func(conn *websocket.Conn) {
			request := &model.JoinRoomRequest{
				RoomID: roomID,
				UserID: userID,
			}
			if err := r.roomUseCase.JoinRoom(conn, request); err != nil {
				return
			}
		})(ctx)
	}

	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": "WebSocket upgrade required",
	})
}

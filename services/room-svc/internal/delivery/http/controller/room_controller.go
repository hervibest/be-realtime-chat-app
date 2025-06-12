package controller

import (
	"be-realtime-chat-app/services/commoner/helper"
	"be-realtime-chat-app/services/room-svc/internal/delivery/http/middleware"
	"be-realtime-chat-app/services/room-svc/internal/helper/logs"
	"be-realtime-chat-app/services/room-svc/internal/model"
	"be-realtime-chat-app/services/room-svc/internal/usecase"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RoomController interface {
	CreateRoom(ctx *fiber.Ctx) error
	GetActiveRooms(ctx *fiber.Ctx) error
	UserDeleteRoom(ctx *fiber.Ctx) error
	UserGetActiveRooms(ctx *fiber.Ctx) error
}

type roomControllerImpl struct {
	roomUseCase usecase.RoomUseCase
	log         logs.Log
}

func NewRoomController(roomUseCase usecase.RoomUseCase, log logs.Log) RoomController {
	return &roomControllerImpl{
		roomUseCase: roomUseCase,
		log:         log,
	}
}

func (c *roomControllerImpl) CreateRoom(ctx *fiber.Ctx) error {
	request := new(model.CreateRoomRequest)
	if err := ctx.BodyParser(request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	response, err := c.roomUseCase.CreateRoom(ctx.Context(), request)
	if err != nil {
		if validatonErrs, ok := err.(*helper.UseCaseValError); ok {
			return helper.ErrValidationResponseJSON(ctx, validatonErrs)
		}
		return helper.ErrUseCaseResponseJSON(ctx, "Create room error : ", err, c.log)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.RoomResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *roomControllerImpl) UserDeleteRoom(ctx *fiber.Ctx) error {
	roomUUID := ctx.Params("roomUUID")
	_, err := uuid.Parse(roomUUID)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid room uuid")
	}

	user := middleware.GetUser(ctx)
	request := &model.UserDeleteRoomRequest{
		RoomUUID: roomUUID,
		UserID:   user.ID,
	}

	if err := c.roomUseCase.UserDeleteRoom(ctx.Context(), request); err != nil {
		if validatonErrs, ok := err.(*helper.UseCaseValError); ok {
			return helper.ErrValidationResponseJSON(ctx, validatonErrs)
		}
		return helper.ErrUseCaseResponseJSON(ctx, "Create room error : ", err, c.log)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[any]{
		Success: true,
	})
}

func (c *roomControllerImpl) GetActiveRooms(ctx *fiber.Ctx) error {
	responses, err := c.roomUseCase.GetActiveRooms(ctx.Context())
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Get active rooms error : ", err, c.log)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*[]*model.RoomResponse]{
		Success: true,
		Data:    responses,
	})
}

func (c *roomControllerImpl) UserGetActiveRooms(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)
	responses, err := c.roomUseCase.UserGetActiveRooms(ctx.Context(), user.ID)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "User get active rooms error : ", err, c.log)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*[]*model.RoomResponse]{
		Success: true,
		Data:    responses,
	})
}

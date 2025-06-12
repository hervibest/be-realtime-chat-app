package controller

import (
	"be-realtime-chat-app/services/chat-query-svc/internal/model"
	"be-realtime-chat-app/services/chat-query-svc/internal/usecase"
	"be-realtime-chat-app/services/commoner/helper"
	"be-realtime-chat-app/services/commoner/logs"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type QueryController interface {
	SearchMessages(ctx *fiber.Ctx) error
}

type queryControllerImpl struct {
	queryUseCase usecase.QueryUseCase
	log          logs.Log
}

func NewQueryController(queryUseCase usecase.QueryUseCase, log logs.Log) QueryController {
	return &queryControllerImpl{
		queryUseCase: queryUseCase,
		log:          log,
	}
}

func (c *queryControllerImpl) SearchMessages(ctx *fiber.Ctx) error {
	request := new(model.SearchParams)
	request.Username = ctx.Query("username")
	request.RoomID = ctx.Params("room_id")
	request.Content = ctx.Query("content")
	request.Limit = ctx.QueryInt("limit", 10)

	response, err := c.queryUseCase.SearchMessages(ctx.Context(), request)
	if err != nil {
		if validatonErrs, ok := err.(*helper.UseCaseValError); ok {
			return helper.ErrValidationResponseJSON(ctx, validatonErrs)
		}
		return helper.ErrUseCaseResponseJSON(ctx, "Search messages error : ", err, c.log)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*[]*model.MessageResponse]{
		Success: true,
		Data:    response,
	})
}

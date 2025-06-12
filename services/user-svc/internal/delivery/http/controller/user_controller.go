package controller

import (
	"be-realtime-chat-app/services/commoner/helper"
	"be-realtime-chat-app/services/commoner/logs"
	"be-realtime-chat-app/services/user-svc/internal/delivery/http/middleware"
	"be-realtime-chat-app/services/user-svc/internal/model"
	"be-realtime-chat-app/services/user-svc/internal/usecase"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserControler interface {
	CurrentUser(ctx *fiber.Ctx) error
	LoginUser(ctx *fiber.Ctx) error
	RegisterUser(ctx *fiber.Ctx) error
	UserLogout(ctx *fiber.Ctx) error
}
type userControlerImpl struct {
	userUC usecase.UserUseCase
	logs   logs.Log
}

func NewUserController(userUC usecase.UserUseCase, logs logs.Log) UserControler {
	return &userControlerImpl{
		userUC: userUC,
		logs:   logs,
	}
}

func (c *userControlerImpl) LoginUser(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	employee, token, err := c.userUC.LoginUser(ctx.UserContext(), request)
	if err != nil {
		if validatonErrs, ok := err.(*helper.UseCaseValError); ok {
			return helper.ErrValidationResponseJSON(ctx, validatonErrs)
		}
		return helper.ErrUseCaseResponseJSON(ctx, "Login employee error : ", err, c.logs)
	}

	response := map[string]interface{}{
		"employee": employee,
		"token":    token,
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
		Data:    response,
	})
}

// TODO STILL RANDOM UUID CREATEDBY
func (c *userControlerImpl) RegisterUser(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	if err := helper.SetBaseModel(ctx, request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid IP address")
	}

	request.CreatedBy = uuid.Must(uuid.NewRandom())
	request.UpdatedBy = request.CreatedBy

	employee, err := c.userUC.RegisterUser(ctx.UserContext(), request)
	if err != nil {
		if validatonErrs, ok := err.(*helper.UseCaseValError); ok {
			return helper.ErrValidationResponseJSON(ctx, validatonErrs)
		}
		return helper.ErrUseCaseResponseJSON(ctx, "Register user error : ", err, c.logs)
	}

	response := map[string]interface{}{
		"employee": employee,
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
		Data:    response,
	})
}

func (c *userControlerImpl) CurrentUser(ctx *fiber.Ctx) error {
	employee := middleware.GetUser(ctx)

	employeeResponse, err := c.userUC.CurrentUser(ctx.Context(), employee.Username)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Current error : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*model.UserResponse]{
		Success: true,
		Data:    employeeResponse,
	})
}

func (c *userControlerImpl) UserLogout(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)

	request := new(model.LogoutUserRequest)
	request.UserId = user.ID
	request.AccessToken = user.Token
	request.ExpiresAt = user.ExpiresAt

	if err := c.userUC.Logout(ctx.Context(), request); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Logout error : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
	})
}

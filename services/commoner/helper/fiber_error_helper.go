package helper

import (
	errorcode "be-realtime-chat-app/services/commoner/constant/errcode"
	"be-realtime-chat-app/services/commoner/logs"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

type BodyParseErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type ValidationErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func ErrCustomResponseJSON(ctx *fiber.Ctx, status int, message string) error {
	return ctx.Status(status).JSON(ErrorResponse{
		Success: false,
		Message: message,
	})
}

func ErrBodyParserResponseJSON(ctx *fiber.Ctx, err error) error {
	return ctx.Status(http.StatusBadRequest).JSON(BodyParseErrorResponse{
		Success: false,
		Message: "Invalid request. Please check the submitted data.",
		Errors:  err.Error(),
	})
}

func ErrValidationResponseJSON(ctx *fiber.Ctx, validatonErrs *UseCaseValError) error {
	return ctx.Status(http.StatusUnprocessableEntity).JSON(ValidationErrorResponse{
		Success: false,
		Message: "Validation error",
		Errors:  validatonErrs.GetValidationErrors(),
	})
}

func ErrUseCaseResponseJSON(ctx *fiber.Ctx, msg string, err error, logs logs.Log) error {
	if appErr, ok := err.(*AppError); ok {
		logs.Info(fmt.Sprintf("UseCase error in controller : %s [%s]: %v", msg, appErr.Code, appErr.Err))
		if appErr.Err != nil {
			logs.Error(fmt.Sprintf("Internal error in controller : %s [%s]: %v", msg, appErr.Code, appErr.Err.Error()))
		} else {
			logs.Info(fmt.Sprintf("Client error in controller : %s [%s]: %v", msg, appErr.Code, appErr.Message))
		}

		return ctx.Status(appErr.HTTPStatus()).JSON(ErrorResponse{
			Success: false,
			Message: appErr.Message,
		})
	}

	return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong. Please try again later")
}

func MultipleULIDSliceParser(ulidSlice []string) error {
	invalidIds := make([]string, 0)
	for _, id := range ulidSlice {
		if _, err := ulid.Parse(id); err != nil {
			invalidIds = append(invalidIds, id)
		}
	}
	if len(invalidIds) != 0 {
		return NewUseCaseError(errorcode.ErrInvalidArgument, fmt.Sprintf("Invalid ids : %s", invalidIds))
	}
	return nil
}

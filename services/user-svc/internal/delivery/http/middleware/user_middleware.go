package middleware

import (
	"be-realtime-chat-app/services/commoner/helper"
	"be-realtime-chat-app/services/commoner/logs"
	"be-realtime-chat-app/services/user-svc/internal/model"
	"be-realtime-chat-app/services/user-svc/internal/usecase"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// TODO SHOULD TOKEN VALIDATED ?
func NewUserAuth(userUC usecase.UserUseCase, validator helper.CustomValidator, logs logs.Log) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token := strings.TrimPrefix(ctx.Get("Authorization", ""), "Bearer ")
		if token == "" || token == "NOT_FOUND" {
			return fiber.NewError(http.StatusUnauthorized, "Unauthorized access")
		}

		user, err := userUC.VerifyUser(ctx.UserContext(), token)
		if err != nil {
			return helper.ErrUseCaseResponseJSON(ctx, "Authenticate user error : ", err, logs)
		}

		ctx.Locals("user", user)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.AuthResponse {
	return ctx.Locals("user").(*model.AuthResponse)
}

package middleware

import (
	"be-realtime-chat-app/services/chat-realtime-svc/internal/adapter"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/helper/logs"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/model"
	"be-realtime-chat-app/services/commoner/helper"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// TODO SHOULD TOKEN VALIDATED ?
func NewUserAuth(userAdapter adapter.UserAdapter, logs logs.Log) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		logs.Info(ctx.Get("Authorization", ""))
		token := strings.TrimPrefix(ctx.Get("Authorization", ""), "Bearer ")
		if token == "" || token == "NOT_FOUND" {
			return fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized access")
		}

		logs.Info(token)

		authResponse, err := userAdapter.AuthenticateUser(ctx.UserContext(), token)
		if err != nil {
			logs.Info(token)
			return helper.ErrUseCaseResponseJSON(ctx, "Authenticate user : ", err, logs)
		}

		logs.Info("User authenticated", zap.String("user_id", authResponse.GetUser().GetId()))

		auth := &model.AuthResponse{
			ID:       authResponse.GetUser().GetId(),
			Username: authResponse.GetUser().GetUsername(),
			Email:    authResponse.GetUser().GetEmail()}

		ctx.Locals("user", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.AuthResponse {
	return ctx.Locals("user").(*model.AuthResponse)
}

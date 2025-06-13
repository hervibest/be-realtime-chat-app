package route

import (
	"be-realtime-chat-app/services/chat-realtime-svc/internal/delivery/http/controller"

	"github.com/gofiber/fiber/v2"
)

type RoomRoute struct {
	app            *fiber.App
	chatController controller.ChatController
	userMiddleware fiber.Handler
}

func NewRoomRoute(app *fiber.App, chatController controller.ChatController, userMiddleware fiber.Handler) *RoomRoute {
	return &RoomRoute{
		app:            app,
		chatController: chatController,
		userMiddleware: userMiddleware,
	}
}

func (r *RoomRoute) RegisterRoutes() {
	chatRoute := r.app.Group("/api/v1/chat", r.userMiddleware)
	chatRoute.Get("/join/:roomID", r.chatController.JoinRoom)
}

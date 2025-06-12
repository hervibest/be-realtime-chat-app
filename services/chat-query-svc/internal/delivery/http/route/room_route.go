package route

import (
	"be-realtime-chat-app/services/chat-query-svc/internal/delivery/http/controller"

	"github.com/gofiber/fiber/v2"
)

type RoomRoute struct {
	app             *fiber.App
	queryController controller.QueryController
	userMiddleware  fiber.Handler
}

func NewRoomRoute(app *fiber.App, queryController controller.QueryController, userMiddleware fiber.Handler) *RoomRoute {
	return &RoomRoute{
		app:             app,
		queryController: queryController,
		userMiddleware:  userMiddleware,
	}
}

func (r *RoomRoute) RegisterRoutes() {
	overtimePeriodRoute := r.app.Group("/api/v1/room", r.userMiddleware)
	overtimePeriodRoute.Get("/:roomID/message", r.queryController.SearchMessages)

}

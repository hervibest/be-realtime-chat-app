package route

import (
	"be-realtime-chat-app/services/room-svc/internal/delivery/http/controller"

	"github.com/gofiber/fiber/v2"
)

type RoomRoute struct {
	app            *fiber.App
	roomController controller.RoomController
	userMiddleware fiber.Handler
}

func NewRoomRoute(app *fiber.App, roomController controller.RoomController, userMiddleware fiber.Handler) *RoomRoute {
	return &RoomRoute{
		app:            app,
		roomController: roomController,
		userMiddleware: userMiddleware,
	}
}

func (r *RoomRoute) RegisterRoutes() {
	overtimePeriodRoute := r.app.Group("/api/v1/room", r.userMiddleware)
	overtimePeriodRoute.Post("/", r.roomController.CreateRoom)
	overtimePeriodRoute.Get("/", r.roomController.GetActiveRooms)
	overtimePeriodRoute.Delete("/:roomUUID", r.roomController.UserDeleteRoom)
	overtimePeriodRoute.Get("/owned", r.roomController.UserGetActiveRooms)

}

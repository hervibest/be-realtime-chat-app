package model

type CreateRoomRequest struct {
	Name   string `json:"name" validate:"required,min=3,max=50"`
	UserID string `json:"user_id" validate:"required"`
}

type UserDeleteRoomRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	RoomUUID string `json:"room_uuid" validate:"required"`
}

type JoinRoomRequest struct {
	RoomID string `json:"room_id" query:"room_id" validate:"required"`
	UserID string `json:"user_id" query:"user_id" validate:"required"`
}

type GetRoomRequest struct {
	RoomID string `json:"room_id" validate:"required"`
}

type RoomResponse struct {
	ID        string `json:"id"`
	UUID      string `json:"uuid,omitempty"`
	Name      string `json:"name"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at,omitempty"`
}

package model

type CreateRoomRequest struct {
	Name   string `json:"name" validate:"required,min=3,max=50"`
	UserID string `json:"user_id" validate:"required,uuid"`
}

type JoinRoomRequest struct {
	RoomID string `json:"room_id" validate:"required,uuid"`
	UserID string `json:"user_id" validate:"required,uuid"`
}

type RoomResponse struct {
	ID        string `json:"id"`
	UUID      string `json:"uuid,omitempty"`
	Name      string `json:"name"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at,omitempty"`
}

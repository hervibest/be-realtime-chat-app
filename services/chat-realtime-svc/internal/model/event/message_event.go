package event

import "be-realtime-chat-app/services/commoner/constant/enum"

type Message struct {
	ID         string              `json:"id"`
	UUID       string              `json:"uuid"`
	RoomID     string              `json:"room_id"`
	RoomStatus enum.RoomStatusEnum `json:"room_status"`
	UserID     string              `json:"user_id"`
	Username   string              `json:"username"`
	Content    string              `json:"content"`
	CreatedAt  string              `json:"created_at"`
	DeletedAt  string              `json:"deleted_at"`
}

func (e *Message) GetID() string {
	return e.ID
}

package event

import "be-realtime-chat-app/services/commoner/constant/enum"

type Message struct {
	ID         string              `db:"id"`
	UUID       string              `db:"uuid"`
	RoomID     string              `db:"room_id"`
	RoomStatus enum.RoomStatusEnum `json:"room_status"`

	UserID    string `db:"user_id"`
	Username  string `db:"username"`
	Content   string `db:"content"`
	CreatedAt string `db:"created_at"`
	DeletedAt string `db:"deleted_at"`
}

func (e *Message) GetID() string {
	return e.ID
}

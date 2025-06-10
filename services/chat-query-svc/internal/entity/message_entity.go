package entity

import "time"

type Message struct {
	ID        string     `json:"id"`
	UUID      string     `json:"uuid"`
	RoomID    string     `json:"room_id"`
	UserID    string     `json:"user_id"`
	Username  string     `json:"username"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

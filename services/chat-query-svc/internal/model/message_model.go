package model

type MessageResponse struct {
	ID        string `json:"id"`
	UUID      string `json:"uuid,omitempty"`
	RoomID    string `json:"room_id"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	DeletedAt string `json:"deleted_at,omitempty"`
}

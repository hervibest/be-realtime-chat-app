package model

type SearchParams struct {
	Username string `json:"username" form:"username" query:"username" `
	Content  string `json:"content" form:"content" query:"content"`
	RoomID   string `json:"room_id" form:"room_id" query:"room_id"`
	Limit    int    `json:"limit" form:"limit" query:"limit"`
}

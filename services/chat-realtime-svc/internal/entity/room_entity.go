package entity

import "time"

type Room struct {
	ID        string     `db:"id"`
	UUID      string     `db:"uuid"`
	Name      string     `db:"name"`
	UserID    string     `db:"user_id"`
	CreatedAt *time.Time `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

package entity

import "time"

type User struct {
	ID        string     `db:"id"`
	UUID      string     `db:"uuid"`
	Email     string     `db:"email"`
	Username  string     `db:"username"`
	Password  string     `db:"password"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

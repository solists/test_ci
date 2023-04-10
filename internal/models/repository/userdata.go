package repository

import "time"

type UserData struct {
	ID        uint64    `db:"id"`
	UserID    int64     `db:"user_id"`
	Allowed   bool      `db:"allowed"`
	ChatID    int64     `db:"chat_id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	UserName  string    `db:"user_name"`
	CreatedAt time.Time `db:"created_at"`
}

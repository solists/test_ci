package repository

import (
	"time"
)

type MessageLog struct {
	ID        uint64    `db:"id"`
	UserID    *int64    `db:"user_id"`
	ChatID    int64     `db:"chat_id"`
	MessageID int       `db:"message_id"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}

type MessageLogResp struct {
	Message string `db:"message"`
}
type MessageLogUserIDResp struct {
	UserID  *int64 `db:"user_id"`
	Message string `db:"message"`
	Allowed *bool  `db:"allowed"`
	ChatID  int64  `db:"chat_id"`
}

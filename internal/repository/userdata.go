package repository

import (
	"context"
	"mymod/internal/models/repository"
)

func (r *Repository) InsertUserData(ctx context.Context, req *repository.UserData) (err error) {
	if _, err = r.db.NamedExecContext(ctx, queryUserDataInsert, req); err != nil {
		return
	}

	return
}

func (r *Repository) UpdateUserDataChatID(ctx context.Context, chatID int64, userID int64) (err error) {
	if _, err = r.db.ExecContext(ctx, queryUserDataUpdate, chatID, userID); err != nil {
		return
	}

	return
}

func (r *Repository) GetUserData(ctx context.Context, userID int64) (user *repository.UserData, err error) {
	user = &repository.UserData{}
	if err = r.db.GetContext(ctx, user, queryUserDataSelect, userID); err != nil {
		return nil, err
	}

	return
}

const queryUserDataInsert = `
INSERT INTO user_data (user_id, chat_id, first_name, last_name, user_name)
VALUES (:user_id, :chat_id, :first_name, :last_name, :user_name)
ON CONFLICT (user_id) DO UPDATE SET allowed=excluded.allowed,
chat_id=excluded.chat_id, first_name=excluded.first_name,
last_name=excluded.last_name, user_name=excluded.user_name
`

const queryUserDataSelect = `
SELECT id, user_id, allowed, chat_id, 
       first_name, last_name, user_name, created_at
FROM user_data WHERE user_id = $1
`

const queryUserDataUpdate = `
UPDATE user_data SET chat_id = $1 WHERE user_id = $2
`

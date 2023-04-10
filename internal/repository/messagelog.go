package repository

import (
	"context"
	"fmt"
	"mymod/internal/models/repository"
	"strings"
)

func (r *Repository) InsertMessageLog(ctx context.Context, req *repository.MessageLog) (err error) {
	if _, err = r.db.NamedExecContext(ctx, queryMessageLog, req); err != nil {
		return
	}

	return
}

func (r *Repository) InsertMessageLogs(ctx context.Context, logs []repository.MessageLog) (err error) {
	if len(logs) == 0 {
		return
	}

	queryIns := `
        INSERT INTO message_log (user_id, chat_id, message_id, message)
        VALUES
    `
	argNum := 4
	var placeholders []string
	var args []interface{}
	for i, log := range logs {
		placeholders = append(placeholders, fmt.Sprintf(" ($%d::bigint, $%d::bigint, $%d::bigint, $%d)",
			i*argNum+1, i*argNum+2, i*argNum+3, i*argNum+4))
		args = append(args, log.UserID, log.ChatID, log.MessageID, log.Message)
	}
	queryIns += strings.Join(placeholders, ", ")

	if _, err = r.db.ExecContext(ctx, queryIns, args...); err != nil {
		return
	}

	return
}

func (r *Repository) SelectLastMessageLogByChatID(
	ctx context.Context,
	chatID uint64,
	limit uint64,
) (messages []repository.MessageLogResp, err error) {
	messages = make([]repository.MessageLogResp, 0)
	if err = r.db.SelectContext(ctx, &messages, querySelectMessageLogByChatID, chatID, limit); err != nil {
		return nil, err
	}

	return
}
func (r *Repository) GetMessageLogWithUserData(
	ctx context.Context,
	chatID int64,
	limit int,
) (resp []repository.MessageLogUserIDResp, err error) {
	resp = make([]repository.MessageLogUserIDResp, 0)
	if err = r.db.SelectContext(ctx, &resp, querySelectMessageLogUserData, chatID, limit); err != nil {
		return nil, err
	}

	return
}

const querySelectMessageLogByChatID = `
SELECT message
FROM message_log
WHERE chat_id = $1
ORDER BY message_id DESC
LIMIT $2;
`

const queryMessageLog = `
INSERT INTO message_log (user_id, chat_id, message_id, message)
VALUES (:user_id, :chat_id, :message_id, :message)
`

const querySelectMessageLogUserData = `
SELECT ml.message, ml.chat_id, ud.user_id, ud.allowed
FROM message_log ml
LEFT JOIN user_data ud ON ml.chat_id = ud.chat_id
WHERE ml.chat_id = $1
ORDER BY ml.message_id DESC
LIMIT $2
`

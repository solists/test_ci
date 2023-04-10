package repository

import (
	"context"
	"mymod/internal/models/repository"

	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=${GOFILE} -destination=mock/mock_${GOFILE}
type IRepository interface {
	AddUsage(ctx context.Context, req *repository.UsageInsert) (err error)
	GetUserData(ctx context.Context, userID int64) (user *repository.UserData, err error)
	InsertUserData(ctx context.Context, req *repository.UserData) (err error)
	UpdateUserDataChatID(ctx context.Context, chatID int64, userID int64) (err error)
	InsertMessageLog(ctx context.Context, req *repository.MessageLog) (err error)
	InsertMessageLogs(ctx context.Context, logs []repository.MessageLog) (err error)
	SelectLastMessageLogByChatID(ctx context.Context, chatID uint64, limit uint64) (messages []repository.MessageLogResp, err error)
	GetMessageLogWithUserData(ctx context.Context, chatID int64, limit int) ([]repository.MessageLogUserIDResp, error)
}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(
	db *sqlx.DB,
) *Repository {
	return &Repository{
		db: db,
	}
}

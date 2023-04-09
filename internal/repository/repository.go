package repository

import (
	"context"
	"mymod/internal/models/repository"

	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=${GOFILE} -destination=mock/mock_${GOFILE}
type IRepository interface {
	AddUsage(ctx context.Context, req *repository.UsageInsert) error
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

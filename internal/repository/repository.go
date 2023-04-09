package repository

import "github.com/jmoiron/sqlx"

//go:generate mockgen -source=${GOFILE} -destination=mock/mock_${GOFILE}
type IRepository interface{}

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
package controller

import (
	"mymod/internal/config"
	"mymod/internal/repository"
)

//go:generate mockgen -source=${GOFILE} -destination=mock/mock_${GOFILE}
type IController interface {
}

type Controller struct {
	repo repository.IRepository
	cfg  *config.Config
}

func NewController(
	repo repository.IRepository,
	cfg *config.Config,
) *Controller {
	return &Controller{
		repo: repo,
		cfg:  cfg,
	}
}

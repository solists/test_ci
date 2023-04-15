package tgservice

import (
	"mymod/internal/controller"
	"mymod/internal/repository"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const defaultLogLimit = 12

var (
	queryRequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "myapp_qet_query_requests_total",
			Help: "Total number of requests received",
		},
		[]string{"user"},
	)
)

type Service struct {
	ctrl       controller.IController
	repo       repository.IRepository
	downloader *Downloader
}

func NewService(
	repo repository.IRepository,
	ctrl controller.IController,
	downloader *Downloader,
) *Service {
	return &Service{
		repo:       repo,
		ctrl:       ctrl,
		downloader: downloader,
	}
}

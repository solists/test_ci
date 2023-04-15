package controller

import (
	"context"
	"mymod/internal/client/openai"
	"mymod/internal/config"
	models "mymod/internal/models/openai"
	"mymod/internal/repository"
	auditservice "mymod/pkg/audit"
)

//go:generate mockgen -source=${GOFILE} -destination=mock/mock_${GOFILE}
type IController interface {
	GetQuery(ctx context.Context, req *models.GetQueryRequest) (*models.GetQueryResponse, error)
	GetTranscription(ctx context.Context, req *models.GetTranscriptionRequest) (*models.GetTranscriptionResponse, error)
}

type Controller struct {
	repo         repository.IRepository
	cfg          *config.Config
	audit        auditservice.Service
	openaiClient openai.Client
}

func NewController(
	repo repository.IRepository,
	cfg *config.Config,
	audit auditservice.Service,
	openaiClient openai.Client,
) *Controller {
	return &Controller{
		repo:         repo,
		cfg:          cfg,
		audit:        audit,
		openaiClient: openaiClient,
	}
}

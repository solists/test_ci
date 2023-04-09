package openai

import (
	"context"
	"mymod/internal/config"
)
import openai "github.com/sashabaranov/go-openai"

//go:generate mockgen -source=${GOFILE} -destination=mock/mock_${GOFILE}
type Client interface {
	GetQueryOPENAI(ctx context.Context, messages []openai.ChatCompletionMessage) (*openai.ChatCompletionResponse, error)
}

type ClientImpl struct {
	cfg    *config.Config
	client *openai.Client
}

func NewClient(
	cfg *config.Config,
) *ClientImpl {
	client := openai.NewClient(cfg.APIKey)
	return &ClientImpl{
		client: client,
		cfg:    cfg,
	}
}

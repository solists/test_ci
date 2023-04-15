package openai

import (
	"context"
	"github.com/pkg/errors"
	"mymod/internal/config"
)
import openai "github.com/sashabaranov/go-openai"

//go:generate mockgen -source=${GOFILE} -destination=mock/mock_${GOFILE}
type Client interface {
	GetQuery(ctx context.Context, messages []openai.ChatCompletionMessage) (*openai.ChatCompletionResponse, error)
	GetTranscription(ctx context.Context, filePath string) (*openai.AudioResponse, error)
}

var ErrTooBigPrompt = errors.New("too big prompt")

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

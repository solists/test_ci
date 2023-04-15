package openai

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/sashabaranov/go-openai"
)

func (c *ClientImpl) GetTranscription(ctx context.Context, filePath string) (*openai.AudioResponse, error) {
	opName := "GetTranscriptionOpenClient"
	sp, _ := opentracing.StartSpanFromContext(ctx, opName)
	defer sp.Finish()

	resp, err := c.client.CreateTranscription(
		ctx,
		openai.AudioRequest{
			Model:    openai.Whisper1,
			FilePath: filePath,
		},
	)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

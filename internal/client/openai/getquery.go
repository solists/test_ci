package openai

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"github.com/solists/test_ci/pkg/logger"
)

func (c *ClientImpl) GetQueryOPENAI(ctx context.Context, messages []openai.ChatCompletionMessage) (*openai.ChatCompletionResponse, error) {
	opName := "GetQueryOpenClient"
	sp, _ := opentracing.StartSpanFromContext(ctx, opName)
	defer sp.Finish()

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) != 1 {
		logger.Warnf("unusual resp, multiple choices: %v", resp)
	} else if len(resp.Choices) == 0 {
		return nil, errors.New("no choices in resp")
	}

	return &resp, nil
}

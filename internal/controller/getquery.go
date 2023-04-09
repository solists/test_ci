package controller

import (
	"context"
	models "mymod/internal/models/openai"
	"mymod/internal/models/repository"
	auditservice "mymod/pkg/audit"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"github.com/solists/test_ci/pkg/logger"
)

func (c *Controller) GetQuery(ctx context.Context, req *models.GetQueryRequest) (*models.GetQueryResponse, error) {
	opName := "GetQuery"
	sp, _ := opentracing.StartSpanFromContext(ctx, opName)
	defer sp.Finish()
	if len(req.Messages) == 0 {
		return nil, errors.New("GetQuery: empty prompt")
	}

	var messages []openai.ChatCompletionMessage
	for _, v := range req.Messages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: v.Message,
		})
	}

	resp, err := c.openaiClient.GetQueryOPENAI(ctx, messages)

	var status uint64 = 200
	if err != nil {
		status = 400
	}
	c.audit.Log(&auditservice.Log{
		UserID:    req.UserID,
		Data:      messages[len(messages)-1],
		Operation: opName,
		Response:  resp,
		Status:    &status,
	})
	if err != nil {
		return nil, err
	}

	if err = c.repo.AddUsage(ctx, &repository.UsageInsert{
		UserID:        req.UserID,
		UsedPrompt:    uint64(resp.Usage.PromptTokens),
		UsedCompleted: uint64(resp.Usage.CompletionTokens),
		UsedTotal:     uint64(resp.Usage.TotalTokens),
	}); err != nil {
		logger.Errorf("AddUsage: %v", err)
	}

	return &models.GetQueryResponse{Result: resp.Choices[0].Message.Content}, nil
}

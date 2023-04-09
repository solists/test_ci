package service

import (
	"context"
	models "mymod/internal/models/openai"

	"github.com/opentracing/opentracing-go"
	v1 "github.com/solists/test_ci/pkg/pb/myapp/v1"
)

func (s *Service) GetQuery(ctx context.Context, req *v1.GetQueryRequest) (*v1.GetQueryResponse, error) {
	opName := "GetQuery"
	sp, _ := opentracing.StartSpanFromContext(ctx, opName)
	defer sp.Finish()

	if err := verifyQueryRequest(req); err != nil {
		return nil, err
	}

	var messages []models.PromptMessage
	for _, v := range req.Messages {
		messages = append(messages, models.PromptMessage{
			Message: v.Message,
		})
	}

	resp, err := s.ctrl.GetQuery(ctx, &models.GetQueryRequest{
		UserID:   req.UserId,
		Messages: messages,
	})
	if err != nil {
		return nil, err
	}

	return &v1.GetQueryResponse{Result: resp.Result}, nil
}

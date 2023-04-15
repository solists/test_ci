package controller

import (
	"context"
	"github.com/opentracing/opentracing-go"
	models "mymod/internal/models/openai"
	auditservice "mymod/pkg/audit"
)

func (c *Controller) GetTranscription(ctx context.Context, req *models.GetTranscriptionRequest) (*models.GetTranscriptionResponse, error) {
	opName := "GetTranscription"
	sp, _ := opentracing.StartSpanFromContext(ctx, opName)
	defer sp.Finish()

	resp, err := c.openaiClient.GetTranscription(ctx, req.FilePath)

	var status uint64 = 200
	if err != nil {
		status = 400
	}
	c.audit.Log(&auditservice.Log{
		UserID:    req.UserID,
		Data:      req.FilePath,
		Operation: opName,
		Response:  resp,
		Status:    &status,
	})
	if err != nil {
		return nil, err
	}

	return &models.GetTranscriptionResponse{Result: resp.Text}, nil
}

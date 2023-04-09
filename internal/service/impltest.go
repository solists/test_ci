package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	v1 "github.com/solists/test_ci/pkg/pb/myapp/v1"
)

func (s *Service) Hello(ctx context.Context, req *v1.HelloRequest) (*v1.HelloResponse, error) {
	opName := "Hello"
	sp, _ := opentracing.StartSpanFromContext(ctx, opName)
	defer sp.Finish()
	return &v1.HelloResponse{Result: "Hello, world!"}, nil
}

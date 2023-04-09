package service

import (
	"context"

	v1 "github.com/solists/test_ci/pkg/pb/myapp/v1"
)

func (s *Service) Add(ctx context.Context, req *v1.AddRequest) (*v1.AddResponse, error) {
	return &v1.AddResponse{Result: req.A + req.B}, nil
}
func (s *Service) Hello(ctx context.Context, req *v1.HelloRequest) (*v1.HelloResponse, error) {
	return &v1.HelloResponse{Result: "Hello, world!"}, nil
}

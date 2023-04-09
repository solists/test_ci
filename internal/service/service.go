package service

import (
	"mymod/internal/controller"

	v1 "github.com/solists/test_ci/pkg/pb/myapp/v1"
)

type Service struct {
	v1.UnimplementedTgServiceServer
	ctrl controller.IController
}

func NewService(
	ctrl controller.IController,
) *Service {
	return &Service{
		ctrl: ctrl,
	}
}

package service

import (
	v1 "github.com/solists/test_ci/pkg/pb/myapp/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func verifyQueryRequest(req *v1.GetQueryRequest) error {
	if req.UserId == 0 {
		return status.Errorf(codes.InvalidArgument, "userID must not be empty")
	}
	if len(req.Messages) == 0 {
		return status.Errorf(codes.InvalidArgument, "messages must not be empty")
	}

	return nil
}

package grpc

import (
	"context"

	authV1 "github.com/AlpacaLabs/protorepo-auth-go/alpacalabs/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Server) SendMFACode(ctx context.Context, request *authV1.SendMFACodeRequest) (*authV1.SendMFACodeResponse, error) {
	// TODO implement
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}

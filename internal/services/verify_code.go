package services

import (
	"context"

	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) VerifyCode(ctx context.Context, request *mfaV1.VerifyCodeRequest) (*mfaV1.VerifyCodeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}

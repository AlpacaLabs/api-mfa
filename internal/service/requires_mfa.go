package service

import (
	"context"

	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Service) RequiresMfa(ctx context.Context, request *mfaV1.RequiresMfaRequest) (*mfaV1.RequiresMfaResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}

package grpc

import (
	"context"

	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
)

func (s Server) RequiresMfa(ctx context.Context, request *mfaV1.RequiresMfaRequest) (*mfaV1.RequiresMfaResponse, error) {
	return s.service.RequiresMfa(ctx, request)
}

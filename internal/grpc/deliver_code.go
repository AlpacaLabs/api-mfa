package grpc

import (
	"context"

	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
)

func (s Server) DeliverCode(ctx context.Context, request *mfaV1.DeliverCodeRequest) (*mfaV1.DeliverCodeResponse, error) {
	return s.service.DeliverCode(ctx, request)
}

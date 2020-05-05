package grpc

import (
	"context"

	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
)

func (s Server) SendCode(ctx context.Context, request *mfaV1.SendCodeRequest) (*mfaV1.SendCodeResponse, error) {
	return s.service.SendCode(ctx, request)
}

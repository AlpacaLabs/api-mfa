package grpc

import (
	"context"

	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
)

func (s Server) VerifyCode(ctx context.Context, request *mfaV1.VerifyCodeRequest) (*mfaV1.VerifyCodeResponse, error) {
	return s.service.VerifyCode(ctx, request)
}

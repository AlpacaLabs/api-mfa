package grpc

import (
	"context"

	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
)

func (s Server) GetDeliveryOptions(ctx context.Context, request *mfaV1.GetDeliveryOptionsRequest) (*mfaV1.GetDeliveryOptionsResponse, error) {
	return s.service.GetDeliveryOptions(ctx, request)
}

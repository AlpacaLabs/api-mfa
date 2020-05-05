package services

import (
	"context"

	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// VerifyCode lets clients complete the MFA flow by sending the code and account ID.
func (s *Service) VerifyCode(ctx context.Context, request *mfaV1.VerifyCodeRequest) (*mfaV1.VerifyCodeResponse, error) {
	//accountID := request.AccountId
	//code := request.Code

	// TODO look up non-stale entity by (accountID, code)
	//   Mark entity as used and stale.
	//   Mark all codes for that account as stale.
	//   Obtain JWT from Auth service
	//   Return JWT in response

	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}

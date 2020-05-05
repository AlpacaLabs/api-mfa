package service

import (
	"context"

	"github.com/AlpacaLabs/mfa/internal/db"

	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// VerifyCode lets clients complete the MFA flow by sending the code and account ID.
func (s *Service) VerifyCode(ctx context.Context, request *mfaV1.VerifyCodeRequest) (*mfaV1.VerifyCodeResponse, error) {
	accountID := request.AccountId
	code := request.Code

	if err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		// Verify the code exists for the given account ID
		c, err := tx.GetCodeByCodeAndAccountID(ctx, code, accountID)
		if err != nil {
			return err
		}

		// Mark the code as used!
		if err := tx.MarkAsUsed(ctx, c.Id); err != nil {
			return err
		}

		// Mark all codes for account as stale
		if err := tx.MarkAllAsStale(ctx, accountID); err != nil {
			return err
		}

		// TODO Obtain JWT from Auth service and return in response

		return nil

	}); err != nil {
		return nil, err
	}

	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}

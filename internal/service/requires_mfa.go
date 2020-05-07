package service

import (
	"context"

	"github.com/AlpacaLabs/api-mfa/internal/db"

	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
)

func (s Service) RequiresMfa(ctx context.Context, request *mfaV1.RequiresMfaRequest) (*mfaV1.RequiresMfaResponse, error) {
	var out *mfaV1.RequiresMfaResponse
	accountID := request.AccountId

	err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		requires, err := tx.RequiresMfa(ctx, accountID)
		if err != nil {
			return err
		}

		out = &mfaV1.RequiresMfaResponse{
			RequiresMfa: requires,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}

package services

import (
	"context"

	"google.golang.org/grpc"

	"github.com/AlpacaLabs/mfa/internal/db"
	hermesV1 "github.com/AlpacaLabs/protorepo-hermes-go/alpacalabs/hermes/v1"
	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
)

func (s *Service) SendCode(ctx context.Context, request *mfaV1.SendCodeRequest) (*mfaV1.SendCodeResponse, error) {

	var out *mfaV1.SendCodeResponse

	err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {

		// TODO use transactional outbox pattern instead of sending SMS before the transaction commits
		hermesConn, err := grpc.Dial(s.config.HermesGRPCAddress)
		if err != nil {
			return err
		}
		smsClient := hermesV1.NewSendSmsServiceClient(hermesConn)

		// Send text
		_, err = smsClient.SendSms(ctx, &hermesV1.SendSmsRequest{
			To:      "555-555-5555",
			Message: "hello from golang",
		})

		if err != nil {
			return err
		}

		out = &mfaV1.SendCodeResponse{}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}

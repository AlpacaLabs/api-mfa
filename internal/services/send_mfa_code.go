package services

import (
	"context"
	"math/rand"
	"time"

	"google.golang.org/grpc"

	clock "github.com/AlpacaLabs/go-timestamp"
	"github.com/AlpacaLabs/mfa/internal/db"
	authV1 "github.com/AlpacaLabs/protorepo-auth-go/alpacalabs/auth/v1"
	hermesV1 "github.com/AlpacaLabs/protorepo-hermes-go/alpacalabs/hermes/v1"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
)

func newCode(accountID string) authV1.MFACode {
	id := xid.New().String()
	now := time.Now()
	return authV1.MFACode{
		Id:        id,
		Code:      randSeq(6),
		CreatedAt: clock.TimeToTimestamp(now),
		ExpiresAt: clock.TimeToTimestamp(now.Add(time.Minute * 30)),
		AccountId: accountID,
	}
}

func (s *Service) SendMFACode(ctx context.Context, in *authV1.SendMFACodeRequest) (*authV1.SendMFACodeResponse, error) {

	accountID := in.AccountId

	var out *authV1.SendMFACodeResponse

	err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {

		// Create an entity
		code := newCode(accountID)
		if err := tx.CreateCode(ctx, code); err != nil {
			return err
		}

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

		out = &authV1.SendMFACodeResponse{}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randSeq(n int) string {
	b := make([]rune, n)
	start := time.Now()
	s1 := rand.NewSource(time.Now().UnixNano())
	log.Println("Seed end time:", time.Since(start))
	r1 := rand.New(s1)
	for i := range b {
		b[i] = letters[r1.Intn(len(letters))]
	}
	return string(b)
}

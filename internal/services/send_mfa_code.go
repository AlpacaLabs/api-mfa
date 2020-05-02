package services

import (
	"context"
	"math/rand"
	"time"

	clock "github.com/AlpacaLabs/go-timestamp"
	"github.com/AlpacaLabs/mfa/internal/db"
	authV1 "github.com/AlpacaLabs/protorepo-auth-go/alpacalabs/auth/v1"
	"github.com/rs/xid"
	"github.com/sfreiberg/gotwilio"
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

		// TODO transactional outbox pattern
		// TODO send phone number with request
		twilioClient := gotwilio.NewTwilioClient(s.config.TwilioAccountSID, s.config.TwilioAuthToken)
		if err := s.SendSms(ctx, SendSmsInput{
			TwilioClient:      twilioClient,
			TwilioPhoneNumber: s.config.TwilioPhoneNumber,
			To:                "555-555-5555",
			Message:           "hello from golang",
		}); err != nil {
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

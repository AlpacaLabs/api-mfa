package services

import (
	"context"
	"math/rand"
	"time"

	"google.golang.org/grpc"

	"github.com/AlpacaLabs/mfa/internal/db"
	hermesV1 "github.com/AlpacaLabs/protorepo-hermes-go/alpacalabs/hermes/v1"
	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
	log "github.com/sirupsen/logrus"
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

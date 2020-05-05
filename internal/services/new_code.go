package services

import (
	"time"

	code "github.com/AlpacaLabs/go-random-code"

	clock "github.com/AlpacaLabs/go-timestamp"
	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
	"github.com/rs/xid"
)

func newCode(accountID string) mfaV1.MFACode {
	id := xid.New().String()
	now := time.Now()
	return mfaV1.MFACode{
		Id:        id,
		Code:      code.New(),
		CreatedAt: clock.TimeToTimestamp(now),
		ExpiresAt: clock.TimeToTimestamp(now.Add(time.Minute * 30)),
		AccountId: accountID,
	}
}

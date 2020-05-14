package entities

import (
	"time"

	clock "github.com/AlpacaLabs/go-timestamp"
	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
)

type MFACode struct {
	ID        string
	AccountID string
	Code      string
	CreatedAt time.Time
	ExpiresAt time.Time
	Used      bool
	Stale     bool
}

func NewMFACodeFromProtobuf(c mfaV1.MFACode) MFACode {
	return MFACode{
		ID:        c.Id,
		AccountID: c.AccountId,
		Code:      c.Code,
		CreatedAt: clock.TimestampToTime(c.CreatedAt),
		ExpiresAt: clock.TimestampToTime(c.ExpiresAt),
		Used:      c.Used,
		Stale:     c.Stale,
	}
}

func (c MFACode) ToProtobuf() *mfaV1.MFACode {
	return &mfaV1.MFACode{
		Id:        c.ID,
		AccountId: c.AccountID,
		Code:      c.Code,
		CreatedAt: clock.TimeToTimestamp(c.CreatedAt),
		ExpiresAt: clock.TimeToTimestamp(c.ExpiresAt),
		Used:      c.Used,
		Stale:     c.Stale,
	}
}

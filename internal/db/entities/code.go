package entities

import (
	clocksql "github.com/AlpacaLabs/go-timestamp-sql"
	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
	"github.com/guregu/null"
)

type MFACode struct {
	ID        string
	AccountID string
	Code      string
	CreatedAt null.Time
	ExpiresAt null.Time
	Used      bool
	Stale     bool
}

func NewMFACodeFromProtobuf(c mfaV1.MFACode) MFACode {
	return MFACode{
		ID:        c.Id,
		AccountID: c.AccountId,
		Code:      c.Code,
		CreatedAt: clocksql.TimestampToNullTime(c.CreatedAt),
		ExpiresAt: clocksql.TimestampToNullTime(c.ExpiresAt),
		Used:      c.Used,
		Stale:     c.Stale,
	}
}

func (c MFACode) ToProtobuf() *mfaV1.MFACode {
	return &mfaV1.MFACode{
		Id:        c.ID,
		AccountId: c.AccountID,
		Code:      c.Code,
		CreatedAt: clocksql.TimestampFromNullTime(c.CreatedAt),
		ExpiresAt: clocksql.TimestampFromNullTime(c.ExpiresAt),
		Used:      c.Used,
		Stale:     c.Stale,
	}
}

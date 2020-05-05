package entities

import (
	"time"

	clock "github.com/AlpacaLabs/go-timestamp"
	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
	"github.com/golang/protobuf/ptypes/timestamp"
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
		CreatedAt: timestampToNullTime(c.CreatedAt),
		ExpiresAt: timestampToNullTime(c.ExpiresAt),
		Used:      c.Used,
		Stale:     c.Stale,
	}
}

func (c MFACode) ToProtobuf() *mfaV1.MFACode {
	return &mfaV1.MFACode{
		Id:        c.ID,
		AccountId: c.AccountID,
		Code:      c.Code,
		CreatedAt: clock.TimeToTimestamp(c.CreatedAt.ValueOrZero()),
		ExpiresAt: clock.TimeToTimestamp(c.ExpiresAt.ValueOrZero()),
		Used:      c.Used,
		Stale:     c.Stale,
	}
}

// TODO add to AlpacaLabs/go-timestamp?
func timestampToNullTime(in *timestamp.Timestamp) null.Time {
	t := clock.TimestampToTime(in)

	var nt null.Time
	if t.IsZero() {
		nt = null.NewTime(time.Time{}, false)
	} else {
		nt = null.TimeFrom(t)
	}
	return nt
}

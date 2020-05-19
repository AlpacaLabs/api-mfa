package db

import (
	"context"
	"fmt"

	hermesV1 "github.com/AlpacaLabs/protorepo-hermes-go/alpacalabs/hermes/v1"
	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
	"github.com/golang/protobuf/proto"

	"github.com/AlpacaLabs/api-mfa/internal/db/entities"
	"github.com/jackc/pgx/v4"
)

const (
	// TableForSendEvent is the name of the transactional outbox (database table)
	// from which we read "jobs" or "events" that need to get sent to a message broker.
	//
	// All records in this table will have the same type of catalyst, e.g.,
	// mfaV1.DeliverCodeRequest.
	//
	// The type of payload protocol buffer may vary, depending on the catalyst.
	// For example, if the catalyst provides an email address ID, the payload
	// will be a hermesV1.SendEmailRequest, whereas if the catalyst provides a
	// phone number ID, the payload will be a hermesV1.SendSmsRequest.
	TableForSendEvent = "txob"
)

type TransactionalOutbox interface {
	ReadEvent(ctx context.Context) (e entities.SendEvent, err error)
	CreateEvent(ctx context.Context, e entities.SendEvent) error

	MarkEventAsSent(ctx context.Context, eventID string) error
}

type outboxImpl struct {
	tx pgx.Tx
}

func (tx *outboxImpl) ReadEvent(ctx context.Context) (e entities.SendEvent, err error) {
	queryTemplate := `
SELECT 
  event_id, trace_id, sampled, sent, catalyst, payload
  FROM %s
  WHERE sent = FALSE
  LIMIT 1
`

	query := fmt.Sprintf(queryTemplate, TableForSendEvent)

	row := tx.tx.QueryRow(ctx, query)

	var catalystBytes []byte
	var payloadBytes []byte

	if err := row.Scan(&e.EventId, &e.TraceId, &e.Sampled, &e.Sent, &catalystBytes, &payloadBytes); err != nil {
		return e, err
	}

	catalyst := &mfaV1.DeliverCodeRequest{}

	if err := proto.Unmarshal(catalystBytes, catalyst); err != nil {
		return e, err
	}

	if catalyst.GetEmailAddressId() != "" {
		payload := &hermesV1.SendEmailRequest{}
		if err := proto.Unmarshal(payloadBytes, payload); err != nil {
			return e, err
		}
	} else if catalyst.GetPhoneNumberId() != "" {
		payload := &hermesV1.SendSmsRequest{}
		if err := proto.Unmarshal(payloadBytes, payload); err != nil {
			return e, err
		}
	} else {
		// TODO return err
	}

	return e, nil
}

func (tx *outboxImpl) CreateEvent(ctx context.Context, e entities.SendEvent) error {
	queryTemplate := `
INSERT INTO %s(
  event_id, trace_id, sampled, sent, catalyst, payload
) 
VALUES($1, $2, $3, $4, $5, $6)
`

	query := fmt.Sprintf(queryTemplate, TableForSendEvent)

	catalystBytes, err := proto.Marshal(&e.Catalyst)
	if err != nil {
		return err
	}

	payloadBytes, err := proto.Marshal(e.Payload)
	if err != nil {
		return err
	}

	if _, err := tx.tx.Exec(ctx, query, e.EventId, e.TraceId, e.Sampled, e.Sent, catalystBytes, payloadBytes); err != nil {
		return err
	}

	return nil
}

func (tx *outboxImpl) MarkEventAsSent(ctx context.Context, eventID string) error {
	queryTemplate := `
UPDATE %s
  SET sent = TRUE
  WHERE event_id = $1
`
	query := fmt.Sprintf(queryTemplate, TableForSendEvent)
	_, err := tx.tx.Exec(ctx, query, eventID)
	return err
}

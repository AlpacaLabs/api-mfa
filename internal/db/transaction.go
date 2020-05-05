package db

import (
	"context"
	"database/sql"

	"github.com/AlpacaLabs/mfa/internal/db/entities"

	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"

	"github.com/golang-sql/sqlexp"
)

type Transaction interface {
	CreateCode(ctx context.Context, code mfaV1.MFACode) error
	GetCode(ctx context.Context, id string) (*mfaV1.MFACode, error)
}

type txImpl struct {
	tx *sql.Tx
}

func (tx *txImpl) CreateCode(ctx context.Context, in mfaV1.MFACode) error {
	var q sqlexp.Querier
	q = tx.tx

	c := entities.NewMFACodeFromProtobuf(in)

	query := `
INSERT INTO authentication_code(
  id, code, created_timestamp, expiration_timestamp, stale, used, account_id
) 
VALUES($1, $2, $3, $4, $5, $6, $7)
`

	_, err := q.ExecContext(ctx, query, c.ID, c.Code, c.CreatedAt, c.ExpiresAt, c.Stale, c.Used, c.AccountID)

	return err
}

func (tx *txImpl) GetCode(ctx context.Context, id string) (*mfaV1.MFACode, error) {
	var q sqlexp.Querier
	q = tx.tx

	var c entities.MFACode

	query := `
SELECT id, code, created_timestamp, expiration_timestamp, stale, used, account_id 
FROM authentication_code
WHERE id=$1
AND stale=FALSE
`

	row := q.QueryRowContext(ctx, query, id)

	err := row.Scan(&c.ID, &c.Code, &c.CreatedAt, &c.ExpiresAt, &c.Stale, &c.Used, &c.AccountID)
	if err != nil {
		return nil, err
	}

	return c.ToProtobuf(), nil
}

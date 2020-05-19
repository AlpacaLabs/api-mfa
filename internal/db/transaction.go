package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/AlpacaLabs/api-mfa/internal/db/entities"

	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
)

const (
	TableForMFACodes = "authentication_code"
)

type Transaction interface {
	TransactionalOutbox

	CreateCode(ctx context.Context, code mfaV1.MFACode) error
	GetCode(ctx context.Context, id string) (*mfaV1.MFACode, error)
	VerifyCode(ctx context.Context, code, accountID string) (*mfaV1.MFACode, error)

	RequiresMfa(ctx context.Context, accountID string) (bool, error)

	MarkAsUsed(ctx context.Context, id string) error
	MarkAllAsStale(ctx context.Context, accountID string) error
}

type txImpl struct {
	tx pgx.Tx
	outboxImpl
}

func newTransaction(tx pgx.Tx) Transaction {
	return &txImpl{
		tx: tx,
		outboxImpl: outboxImpl{
			tx: tx,
		},
	}
}

func (tx *txImpl) CreateCode(ctx context.Context, in mfaV1.MFACode) error {
	c := entities.NewMFACodeFromProtobuf(in)

	queryTemplate := `
INSERT INTO %s(
  id, code, created_at, expires_at, stale, used, account_id
) 
VALUES($1, $2, $3, $4, $5, $6, $7)
`

	query := fmt.Sprintf(queryTemplate, TableForMFACodes)
	_, err := tx.tx.Exec(ctx, query, c.ID, c.Code, c.CreatedAt, c.ExpiresAt, c.Stale, c.Used, c.AccountID)

	return err
}

func (tx *txImpl) GetCode(ctx context.Context, id string) (*mfaV1.MFACode, error) {
	var c entities.MFACode

	queryTemplate := `
SELECT id, code, created_at, expires_at, stale, used, account_id 
FROM %s
WHERE id=$1
AND stale=FALSE
`

	query := fmt.Sprintf(queryTemplate, TableForMFACodes)
	row := tx.tx.QueryRow(ctx, query, id)

	err := row.Scan(&c.ID, &c.Code, &c.CreatedAt, &c.ExpiresAt, &c.Stale, &c.Used, &c.AccountID)
	if err != nil {
		return nil, err
	}

	return c.ToProtobuf(), nil
}

func (tx *txImpl) VerifyCode(ctx context.Context, code, accountID string) (*mfaV1.MFACode, error) {
	var c entities.MFACode

	queryTemplate := `
SELECT id, code, created_at, expires_at, stale, used, account_id 
FROM %s
WHERE code=$1
AND account_id=$2
AND stale=FALSE
`

	query := fmt.Sprintf(queryTemplate, TableForMFACodes)
	row := tx.tx.QueryRow(ctx, query, code, accountID)

	err := row.Scan(&c.ID, &c.Code, &c.CreatedAt, &c.ExpiresAt, &c.Stale, &c.Used, &c.AccountID)
	if err != nil {
		return nil, err
	}

	return c.ToProtobuf(), nil
}

func (tx *txImpl) RequiresMfa(ctx context.Context, accountID string) (bool, error) {
	var requiresMFA bool

	queryTemplate := `
SELECT requires_mfa 
FROM account
WHERE id=$1
`

	query := fmt.Sprintf(queryTemplate, TableForMFACodes)
	row := tx.tx.QueryRow(ctx, query, accountID)

	err := row.Scan(&requiresMFA)
	if err != nil {
		return false, err
	}

	return requiresMFA, nil
}

func (tx *txImpl) MarkAsUsed(ctx context.Context, id string) error {
	queryTemplate := `
UPDATE %s SET used = TRUE WHERE id=$1
`

	query := fmt.Sprintf(queryTemplate, TableForMFACodes)
	_, err := tx.tx.Exec(ctx, query, id)

	return err
}

func (tx *txImpl) MarkAllAsStale(ctx context.Context, accountID string) error {
	queryTemplate := `
UPDATE %s SET used = TRUE, stale = TRUE WHERE account_id=$1
`

	query := fmt.Sprintf(queryTemplate, TableForMFACodes)
	_, err := tx.tx.Exec(ctx, query, accountID)

	return err
}

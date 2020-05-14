package db

import (
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/AlpacaLabs/api-mfa/internal/db/entities"

	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
)

type Transaction interface {
	CreateCode(ctx context.Context, code mfaV1.MFACode) error
	GetCode(ctx context.Context, id string) (*mfaV1.MFACode, error)
	GetCodeByCodeAndAccountID(ctx context.Context, code, accountID string) (*mfaV1.MFACode, error)

	CreateTxobForCode(ctx context.Context, in mfaV1.DeliverCodeRequest) error
	RequiresMfa(ctx context.Context, accountID string) (bool, error)

	MarkAsUsed(ctx context.Context, id string) error
	MarkAllAsStale(ctx context.Context, accountID string) error
}

type txImpl struct {
	tx pgx.Tx
}

func (tx *txImpl) CreateCode(ctx context.Context, in mfaV1.MFACode) error {
	c := entities.NewMFACodeFromProtobuf(in)

	query := `
INSERT INTO authentication_code(
  id, code, created_at, expires_at, stale, used, account_id
) 
VALUES($1, $2, $3, $4, $5, $6, $7)
`

	_, err := tx.tx.Exec(ctx, query, c.ID, c.Code, c.CreatedAt, c.ExpiresAt, c.Stale, c.Used, c.AccountID)

	return err
}

func (tx *txImpl) GetCode(ctx context.Context, id string) (*mfaV1.MFACode, error) {
	var c entities.MFACode

	query := `
SELECT id, code, created_at, expires_at, stale, used, account_id 
FROM authentication_code
WHERE id=$1
AND stale=FALSE
`

	row := tx.tx.QueryRow(ctx, query, id)

	err := row.Scan(&c.ID, &c.Code, &c.CreatedAt, &c.ExpiresAt, &c.Stale, &c.Used, &c.AccountID)
	if err != nil {
		return nil, err
	}

	return c.ToProtobuf(), nil
}

func (tx *txImpl) GetCodeByCodeAndAccountID(ctx context.Context, code, accountID string) (*mfaV1.MFACode, error) {
	var c entities.MFACode

	query := `
SELECT id, code, created_at, expires_at, stale, used, account_id 
FROM authentication_code
WHERE code=$1
AND account_id=$2
AND stale=FALSE
`

	row := tx.tx.QueryRow(ctx, query, code, accountID)

	err := row.Scan(&c.ID, &c.Code, &c.CreatedAt, &c.ExpiresAt, &c.Stale, &c.Used, &c.AccountID)
	if err != nil {
		return nil, err
	}

	return c.ToProtobuf(), nil
}

func (tx *txImpl) CreateTxobForCode(ctx context.Context, in mfaV1.DeliverCodeRequest) error {
	query := `
INSERT INTO mfa_code_txob(code_id, sent, email_address_id, phone_number_id) 
 VALUES($1, FALSE, $2, $3)
`

	_, err := tx.tx.Exec(ctx, query, in.CodeId, in.GetEmailAddressId(), in.GetPhoneNumberId())

	return err
}

func (tx *txImpl) RequiresMfa(ctx context.Context, accountID string) (bool, error) {
	var requiresMFA bool

	query := `
SELECT requires_mfa 
FROM account
WHERE id=$1
`

	row := tx.tx.QueryRow(ctx, query, accountID)

	err := row.Scan(&requiresMFA)
	if err != nil {
		return false, err
	}

	return requiresMFA, nil
}

func (tx *txImpl) MarkAsUsed(ctx context.Context, id string) error {
	query := `
UPDATE authentication_code SET used = TRUE WHERE id=$1
`

	_, err := tx.tx.Exec(ctx, query, id)

	return err
}

func (tx *txImpl) MarkAllAsStale(ctx context.Context, accountID string) error {
	query := `
UPDATE authentication_code SET used = TRUE, stale = TRUE WHERE account_id=$1
`

	_, err := tx.tx.Exec(ctx, query, accountID)

	return err
}

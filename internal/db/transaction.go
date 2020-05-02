package db

import (
	"context"
	"database/sql"

	authV1 "github.com/AlpacaLabs/protorepo-auth-go/alpacalabs/auth/v1"
	"github.com/golang-sql/sqlexp"
)

type Transaction interface {
	CreateCode(ctx context.Context, code authV1.MFACode) error
}

type txImpl struct {
	tx *sql.Tx
}

func (tx *txImpl) CreateCode(ctx context.Context, c authV1.MFACode) error {
	var q sqlexp.Querier
	q = tx.tx

	_, err := q.ExecContext(
		ctx,
		"INSERT INTO authentication_code(id, code, created_timestamp, expiration_timestamp, stale, used, account_id) VALUES($1, $2, $3, $4, $5, $6, $7)",
		c.Id, c.Code, c.CreatedAt.Seconds, c.ExpiresAt, c.Stale, c.Used, c.AccountId)

	return err
}

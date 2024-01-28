// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: wallet.sql

package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createWallet = `-- name: CreateWallet :one
INSERT INTO wallets (
    balance
) VALUES (
    $1
)
RETURNING id, balance
`

func (q *Queries) CreateWallet(ctx context.Context, balance pgtype.Numeric) (Wallet, error) {
	row := q.db.QueryRow(ctx, createWallet, balance)
	var i Wallet
	err := row.Scan(&i.ID, &i.Balance)
	return i, err
}

const getWallet = `-- name: GetWallet :one
SELECT id, balance FROM wallets
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetWallet(ctx context.Context, id uuid.UUID) (Wallet, error) {
	row := q.db.QueryRow(ctx, getWallet, id)
	var i Wallet
	err := row.Scan(&i.ID, &i.Balance)
	return i, err
}

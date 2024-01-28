-- name: CreateTransaction :one
INSERT INTO transactions (wallet_from, wallet_to, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListTransactionsByWalletId :many
SELECT *
FROM transactions
WHERE (wallet_from = sqlc.arg(wallet_id) OR wallet_to = sqlc.arg(wallet_id))
ORDER BY time DESC;
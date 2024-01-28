-- name: CreateTransaction :one
INSERT INTO transactions (wallet_from_id, wallet_to_id, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListTransactionsByWalletId :many
SELECT *
FROM transactions
WHERE (wallet_from_id = sqlc.arg(wallet_id) OR wallet_to_id = sqlc.arg(wallet_id))
ORDER BY time DESC;
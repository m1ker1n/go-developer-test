-- name: CreateWallet :one
INSERT INTO wallets (
    balance
) VALUES (
    $1
)
RETURNING *;

-- name: GetWallet :one
SELECT * FROM wallets
WHERE id = $1 LIMIT 1;
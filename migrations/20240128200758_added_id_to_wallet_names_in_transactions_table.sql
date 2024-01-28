-- +goose Up
-- +goose StatementBegin
ALTER TABLE transactions
    RENAME COLUMN wallet_from TO wallet_from_id;
ALTER TABLE transactions
    RENAME COLUMN wallet_to TO wallet_to_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE transactions
    RENAME COLUMN wallet_from_id TO wallet_from;
ALTER TABLE transactions
    RENAME COLUMN wallet_to_id TO wallet_to;
-- +goose StatementEnd

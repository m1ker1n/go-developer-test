-- +goose Up
-- +goose StatementBegin
ALTER TABLE wallets
    ALTER COLUMN balance SET NOT NULL;
ALTER TABLE transactions
    ALTER COLUMN amount SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE wallets
    ALTER COLUMN balance DROP NOT NULL;
ALTER TABLE transactions
    ALTER COLUMN amount DROP NOT NULL;
-- +goose StatementEnd

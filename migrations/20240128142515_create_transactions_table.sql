-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions
(
    id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    time    timestamptz NOT NULL DEFAULT now(),
    wallet_from UUID NOT NULL,
    wallet_to UUID NOT NULL,
    amount decimal
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;
-- +goose StatementEnd
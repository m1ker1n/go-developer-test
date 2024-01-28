-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallets
(
    id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    balance decimal
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE wallets;
-- +goose StatementEnd
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS currency (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    code VARCHAR(3) NOT NULL UNIQUE,
    rate NUMERIC(15, 2) NOT NULL DEFAULT 0.00
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS currency;
-- +goose StatementEnd

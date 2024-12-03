-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wallet (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE, -- Ссылка на пользователя (1:1 связь)
    balance_eur NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    balance_usd NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    balance_rub NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE
);
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON wallet
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_timestamp ON wallet;
DROP FUNCTION IF EXISTS trigger_set_timestamp;
DROP TABLE IF EXISTS wallet;
-- +goose StatementEnd

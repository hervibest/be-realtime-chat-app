-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id CHAR(26) PRIMARY KEY NOT NULL,
    uuid UUID DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE,
    username VARCHAR(255) UNIQUE,
    password VARCHAR(60),
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    deleted_at TIMESTAMPTZ
)

CREATE INDEX idx_users_uuid (uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_uuid;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS rooms (
    id CHAR(26) PRIMARY KEY NOT NULL,
    uuid UUID DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    user_id CHAR(26) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_rooms_uuid ON rooms (uuid);
CREATE INDEX idx_rooms_user_id ON rooms (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS rooms;
DROP INDEX IF EXISTS idx_rooms_uuid;
DROP INDEX IF EXISTS idx_rooms_user_id;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS messaging_service.messages (
    id text,
    uuid uuid,
    room_id text,
    user_id text,
    username text,
    content text,
    created_at timestamp,
    deleted_at timestamp,
    PRIMARY KEY ((room_id), created_at, id)
) WITH CLUSTERING ORDER BY (created_at DESC);

CREATE INDEX IF NOT EXISTS ON messaging_service.messages (id);

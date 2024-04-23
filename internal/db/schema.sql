CREATE TABLE note (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    expires_at_timezone TEXT NOT NULL,
    key_hash BYTEA NOT NULL
);

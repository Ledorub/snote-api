CREATE TABLE note (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    content VARCHAR(1048576) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    expires_at_timezone TEXT NOT NULL,
    key_hash BYTEA NOT NULL

    CONSTRAINT chk_key_hash_32_bytes CHECK (length(key_hash) = 32)
);

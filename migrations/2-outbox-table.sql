CREATE TABLE IF NOT EXISTS outbox
(
    id         BIGSERIAL PRIMARY KEY,
    message_id text    NOT NULL REFERENCES messages (id) ON DELETE CASCADE,
    payload    JSONB   NOT NULL,
    sent       BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_outbox_sent_created_at ON outbox (sent, created_at);
CREATE INDEX IF NOT EXISTS idx_outbox_message_id ON outbox (message_id); 
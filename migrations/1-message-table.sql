CREATE TABLE IF NOT EXISTS messages
(
    id           BIGSERIAL PRIMARY KEY,
    content      VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20)  NOT NULL,
    sent         BOOLEAN      NOT NULL DEFAULT false,
    created_at   TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP             DEFAULT CURRENT_TIMESTAMP
);

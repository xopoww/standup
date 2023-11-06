CREATE TABLE messages (
    id              VARCHAR(16) PRIMARY KEY,
    owner_id        VARCHAR(32) NOT NULL,
    content         TEXT NOT NULL,
    created_at      TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS types (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    name       TEXT        NOT NULL UNIQUE,
    version    INTEGER     NOT NULL DEFAULT 1
);
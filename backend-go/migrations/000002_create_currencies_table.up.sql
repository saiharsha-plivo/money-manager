CREATE TABLE IF NOT EXISTS currencies
(
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    name       TEXT        NOT NULL UNIQUE,
    rate       DECIMAL     NOT NULL,
    version    INTEGER     NOT NULL DEFAULT 1
);

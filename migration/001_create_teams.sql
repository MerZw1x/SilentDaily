-- +goose Up
CREATE SCHEMA IF NOT EXISTS silentdaily;

CREATE TABLE IF NOT EXISTS silentdaily.teams (
    id         BIGSERIAL   PRIMARY KEY,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS silentdaily.teams;

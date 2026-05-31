-- +goose Up
CREATE TABLE IF NOT EXISTS silentdaily.ai_api (
    hash     TEXT    PRIMARY KEY,
    requests INTEGER NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE IF EXISTS silentdaily.ai_api;

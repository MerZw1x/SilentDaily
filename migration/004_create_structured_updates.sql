-- +goose Up
CREATE TABLE IF NOT EXISTS silentdaily.structured_updates (
    id              BIGSERIAL   PRIMARY KEY,
    daily_update_id BIGINT      NOT NULL REFERENCES silentdaily.daily_updates(id),
    progress        TEXT[]      NOT NULL DEFAULT '{}',
    plans           TEXT[]      NOT NULL DEFAULT '{}',
    blockers        TEXT[]      NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS silentdaily.structured_updates;

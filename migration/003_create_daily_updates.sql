-- +goose Up
CREATE TABLE IF NOT EXISTS silentdaily.daily_updates (
    id         BIGSERIAL   PRIMARY KEY,
    member_id  BIGINT      NOT NULL REFERENCES silentdaily.members(id),
    team_id    BIGINT      NOT NULL REFERENCES silentdaily.teams(id),
    raw_text   TEXT        NOT NULL,
    status     TEXT        NOT NULL DEFAULT 'queued',
    attempts   INTEGER     NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS silentdaily.daily_updates;

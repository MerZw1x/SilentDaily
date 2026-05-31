-- +goose Up
CREATE TABLE IF NOT EXISTS silentdaily.digests (
    id          BIGSERIAL   PRIMARY KEY,
    team_id     BIGINT      NOT NULL REFERENCES silentdaily.teams(id),
    date        DATE        NOT NULL,
    lead_digest TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (team_id, date)
);

-- +goose Down
DROP TABLE IF EXISTS silentdaily.digests;

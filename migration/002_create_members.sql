-- +goose Up
CREATE TABLE IF NOT EXISTS silentdaily.members (
    id               BIGSERIAL   PRIMARY KEY,
    team_id          BIGINT      NOT NULL REFERENCES silentdaily.teams(id),
    telegram_user_id BIGINT      NOT NULL UNIQUE,
    name             TEXT        NOT NULL,
    is_lead          BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS silentdaily.members;

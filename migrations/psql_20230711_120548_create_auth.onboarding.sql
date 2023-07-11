
-- +migrate Up
CREATE TABLE IF NOT EXISTS auth.onboarding(
    id uuid NOT NULL PRIMARY KEY,
    client_id BIGINT  NOT NULL,
    request_id VARCHAR (255) NOT NULL,
    user_id VARCHAR (50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS auth.onboarding;

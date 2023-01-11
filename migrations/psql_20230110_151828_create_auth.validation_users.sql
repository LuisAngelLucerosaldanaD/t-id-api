
-- +migrate Up
CREATE TABLE IF NOT EXISTS auth.validation_users(
    id uuid NOT NULL PRIMARY KEY,
    transaction_id VARCHAR (100) NOT NULL,
    user_id VARCHAR (100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS auth.validation_users;

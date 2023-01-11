
-- +migrate Up
CREATE TABLE IF NOT EXISTS auth.users_rol(
    id uuid NOT NULL PRIMARY KEY,
    user_id VARCHAR (100) NOT NULL,
    role_id VARCHAR (100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS auth.users_rol;

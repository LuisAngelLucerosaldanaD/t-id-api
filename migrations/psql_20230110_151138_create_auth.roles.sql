
-- +migrate Up
CREATE TABLE IF NOT EXISTS auth.roles(
    id uuid NOT NULL PRIMARY KEY,
    name VARCHAR (50) NOT NULL,
    description VARCHAR (500) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS auth.roles;


-- +migrate Up
CREATE TABLE IF NOT EXISTS auth.role(
    id uuid NOT NULL PRIMARY KEY,
    name VARCHAR (100) NOT NULL,
    description VARCHAR (255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS auth.role;

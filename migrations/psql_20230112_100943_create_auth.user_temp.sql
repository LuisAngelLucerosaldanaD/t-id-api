
-- +migrate Up
CREATE TABLE IF NOT EXISTS auth.user_temp(
    id uuid NOT NULL PRIMARY KEY,
    full_name VARCHAR (150) NOT NULL,
    surname VARCHAR (100) NOT NULL,
    name VARCHAR (100) NOT NULL,
    picture VARCHAR (150) NOT NULL,
    email VARCHAR (150) NOT NULL,
    domain VARCHAR (150) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS auth.user_temp;

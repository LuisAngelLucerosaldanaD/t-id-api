
-- +migrate Up
CREATE TABLE IF NOT EXISTS cfg.files(
    id BIGSERIAL  NOT NULL PRIMARY KEY,
    path VARCHAR (100) NOT NULL,
    name VARCHAR (100) NOT NULL,
    type CHANGE-THIS-TYPE  NOT NULL,
    user_id VARCHAR (100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS cfg.files;
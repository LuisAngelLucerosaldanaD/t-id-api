
-- +migrate Up
CREATE TABLE IF NOT EXISTS wf.work_validation(
    id BIGSERIAL  NOT NULL PRIMARY KEY,
    status VARCHAR (50) NOT NULL,
    user_id VARCHAR (100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS wf.work_validation;

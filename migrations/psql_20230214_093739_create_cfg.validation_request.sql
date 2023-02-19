
-- +migrate Up
CREATE TABLE IF NOT EXISTS cfg.validation_request(
    id BIGSERIAL  NOT NULL PRIMARY KEY,
    client_id BIGINT  NOT NULL,
    max_num_validation INTEGER  NOT NULL,
    request_id VARCHAR (100) NOT NULL,
    expired_at TIMESTAMP  NOT NULL,
    user_identification VARCHAR (100) NOT NULL,
    status varchar (15) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS cfg.validation_request;

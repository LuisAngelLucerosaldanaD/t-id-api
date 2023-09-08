-- +migrate Up
CREATE TABLE IF NOT EXISTS cfg.client
(
    id           BIGSERIAL    NOT NULL PRIMARY KEY,
    full_name    VARCHAR(150) NOT NULL,
    nit          VARCHAR(50)  NOT NULL,
    banner       VARCHAR(50)  NOT NULL,
    logo_small   VARCHAR(50)  NOT NULL,
    main_color   VARCHAR(50)  NOT NULL,
    second_color VARCHAR(50)  NOT NULL,
    url_redirect VARCHAR(500) NOT NULL,
    url_api      VARCHAR(500) NOT NULL,
    created_at   TIMESTAMP    NOT NULL DEFAULT now(),
    updated_at   TIMESTAMP    NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS cfg.client;

-- +migrate Up
CREATE TABLE IF NOT EXISTS auth.user
(
    id                   uuid         NOT NULL PRIMARY KEY,
    nickname             VARCHAR(50)  NOT NULL UNIQUE,
    email                VARCHAR(50)  NOT NULL UNIQUE,
    password             varchar(150) NOT NULL,
    first_name           varchar(50)  NULL,
    second_name          VARCHAR(50)  NULL,
    first_surname        VARCHAR(50)  NULL,
    second_surname       VARCHAR(50)  NULL,
    age                  int4         NULL,
    type_document        VARCHAR(50)  NULL,
    document_number      VARCHAR(20)  NOT NULL UNIQUE,
    cellphone            VARCHAR(15)  NULL,
    gender               VARCHAR(10)  NULL,
    nationality          VARCHAR(50)  NULL,
    country              varchar(50)  NULL,
    department           VARCHAR(50)  NULL,
    city                 VARCHAR(50)  NULL,
    real_ip              VARCHAR(20)  NOT NULL,
    status_id            int4         NOT NULL,
    failed_attempts      int4         NOT NULL DEFAULT 0,
    block_date           TIMESTAMP    NULL,
    disabled_date        TIMESTAMP    NULL,
    last_login           TIMESTAMP    NULL,
    last_change_password TIMESTAMP    NULL,
    birth_date           TIMESTAMP    NULL,
    verified_code        varchar(150) NULL,
    is_deleted           bool         NULL     default false,
    deleted_at           timestamp    NULL,
    created_at           TIMESTAMP    NOT NULL DEFAULT now(),
    updated_at           TIMESTAMP    NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS auth.user;

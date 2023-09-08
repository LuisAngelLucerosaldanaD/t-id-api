-- +migrate Up
CREATE TABLE IF NOT EXISTS trx.traceability
(
    id          BIGSERIAL    NOT NULL PRIMARY KEY,
    action      VARCHAR(100) NOT NULL,
    type        VARCHAR(50)  NOT NULL,
    description VARCHAR(500) NOT NULL,
    user_id     VARCHAR(100) NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP    NOT NULL DEFAULT now(),
    constraint FK_user_trx foreign key (user_id) references auth.user (id)
);

-- +migrate Down
DROP TABLE IF EXISTS trx.traceability;

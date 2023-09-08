-- +migrate Up
CREATE TABLE IF NOT EXISTS auth.life_test
(
    id           BIGSERIAL    NOT NULL PRIMARY KEY,
    client_id    BIGINT       NOT NULL,
    max_num_test INT          NOT NULL,
    request_id   VARCHAR(100) NOT NULL,
    expired_at   TIMESTAMP    NOT NULL,
    user_id      uuid         NOT NULL,
    status       varchar(15)  NOT NULL,
    created_at   TIMESTAMP    NOT NULL DEFAULT now(),
    updated_at   TIMESTAMP    NOT NULL DEFAULT now(),
    constraint FK_user_life_test foreign key (user_id) references auth.user (id),
    constraint FK_client_life_test foreign key (client_id) references auth.client (id)
);

-- +migrate Down
DROP TABLE IF EXISTS auth.life_test;

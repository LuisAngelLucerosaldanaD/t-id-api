-- +migrate Up
CREATE TABLE IF NOT EXISTS auth.onboarding
(
    id             uuid         NOT NULL PRIMARY KEY,
    client_id      BIGINT       NOT NULL,
    request_id     VARCHAR(100) NOT NULL,
    user_id        uuid         NOT NULL,
    status         varchar(50)  not null,
    transaction_id varchar(100) NULL,
    created_at     TIMESTAMP    NOT NULL DEFAULT now(),
    updated_at     TIMESTAMP    NOT NULL DEFAULT now(),
    constraint FK_user_onboarding foreign key (user_id) references auth.user (id),
    constraint FK_client_onboarding foreign key (client_id) references cfg.client (id)
);

-- +migrate Down
DROP TABLE IF EXISTS auth.onboarding;

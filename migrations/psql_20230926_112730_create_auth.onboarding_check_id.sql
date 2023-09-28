-- +migrate Up
CREATE TABLE IF NOT EXISTS auth.onboarding_check_id
(
    id         BIGSERIAL    NOT NULL PRIMARY KEY,
    user_id    uuid NOT NULL unique,
    ip         VARCHAR(20)  NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT now(),
    updated_at TIMESTAMP    NOT NULL DEFAULT now(),
    constraint fk_user_check_id foreign key (user_id) references auth."user" (id)
);

-- +migrate Down
DROP TABLE IF EXISTS auth.onboarding_check_id;

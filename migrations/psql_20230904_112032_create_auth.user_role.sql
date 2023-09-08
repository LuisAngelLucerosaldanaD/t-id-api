
-- +migrate Up
CREATE TABLE IF NOT EXISTS auth.user_role(
    id uuid NOT NULL PRIMARY KEY,
    user_id uuid NOT NULL,
    role_id uuid NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    constraint fk_role foreign key(role_id) references auth.role(id),
    constraint fk_user foreign key(user_id) references auth.user(id)
);

-- +migrate Down
DROP TABLE IF EXISTS auth.user_role;

-- +migrate Up
CREATE TABLE IF NOT EXISTS cfg.file
(
    id         BIGSERIAL    NOT NULL PRIMARY KEY,
    path       VARCHAR(100) NOT NULL,
    name       VARCHAR(100) NOT NULL,
    type       int4         NOT NULL,
    user_id    uuid         NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT now(),
    updated_at TIMESTAMP    NOT NULL DEFAULT now(),
    CONSTRAINT FK_user_file FOREIGN KEY (user_id) references auth.user (id)
);

-- +migrate Down
DROP TABLE IF EXISTS cfg.file;

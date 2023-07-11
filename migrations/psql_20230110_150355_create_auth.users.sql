
-- +migrate Up
CREATE TABLE IF NOT EXISTS auth.users(
    id uuid NOT NULL PRIMARY KEY,
    type_document VARCHAR (50) NULL,
    document_number BIGINT NULL,
    expedition_date TIMESTAMP NULL,
    email varchar (50)  NOT NULL UNIQUE,
    first_name VARCHAR (50) NULL,
    second_name VARCHAR (50) ,
    second_surname VARCHAR (50) NULL,
    age int4  NULL,
    gender VARCHAR (20) NULL,
    nationality VARCHAR (50) NULL,
    civil_status VARCHAR (50) NULL,
    first_surname VARCHAR (50) NULL,
    birth_date TIMESTAMP NULL,
    country VARCHAR (100) NULL,
    department VARCHAR (100) NULL,
    cellphone VARCHAR (15) NULL,
    city VARCHAR (100) NULL,
    real_ip VARCHAR (50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS auth.users;

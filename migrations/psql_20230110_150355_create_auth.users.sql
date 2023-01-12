
-- +migrate Up
CREATE TABLE IF NOT EXISTS auth.users(
    id uuid NOT NULL PRIMARY KEY,
    type_document VARCHAR (50) NOT NULL,
    document_number BIGINT  NOT NULL UNIQUE,
    expedition_date TIMESTAMP  NOT NULL,
    email varchar (50)  NOT NULL UNIQUE,
    first_name VARCHAR (50) NOT NULL,
    second_name VARCHAR (50) ,
    second_surname VARCHAR (50) NOT NULL,
    age int4  NOT NULL,
    gender VARCHAR (20) NOT NULL,
    nationality VARCHAR (50) NOT NULL,
    civil_status VARCHAR (50) NOT NULL,
    first_surname VARCHAR (50) NOT NULL,
    birth_date TIMESTAMP  NOT NULL,
    country VARCHAR (100) NOT NULL,
    department VARCHAR (100) NOT NULL,
    city VARCHAR (100) NOT NULL,
    real_ip VARCHAR (50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS auth.users;

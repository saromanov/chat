-- +goose Up
BEGIN TRANSACTION;

CREATE TABLE users (
    id int NOT NULL
    email text NOT NULL,
    password text NOT NULL,
    first_name text,
    last_name text,
    created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY(id)
);

CREATE TABLE messages (
    id int NOT NULL
    body text,
     created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY(id)
);

COMMIT TRANSACTION;

-- +goose Down
BEGIN TRANSACTION

DROP TABLE users;
DROp TABLE mesages;
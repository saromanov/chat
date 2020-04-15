-- +goose Up
CREATE TABLE users (
    id int NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    first_name text,
    last_name text,
    created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY(id)
);

CREATE TABLE messages (
    id int NOT NULL,
    body text,
     created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY(id)
);

-- +goose Down


DROP TABLE users;
DROP TABLE messages;

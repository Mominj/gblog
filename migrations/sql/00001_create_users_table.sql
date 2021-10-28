-- +goose Up
-- SQL in this section is executed when the migration is applied.


CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE users;
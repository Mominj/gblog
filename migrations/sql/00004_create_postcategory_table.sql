-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE postcategory
(
    id SERIAL PRIMARY KEY,
    blogid INT NOT NULL,
    catid INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE postcategory;
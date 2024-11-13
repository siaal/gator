-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY key,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    name VARCHAR(100) NOT NULL,
    CONSTRAINT unique_username UNIQUE (name)
);

-- +goose Down
DROP TABLE users;

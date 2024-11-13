-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY key,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name VARCHAR(100) NOT NULL,
    CONSTRAINT unique_username UNIQUE (name)
);

-- +goose Down
DROP TABLE users;

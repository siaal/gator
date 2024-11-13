-- +goose Up
CREATE TABLE feeds (
    id uuid PRIMARY key,
    name text NOT NULL,
    url text NOT NULL,
    user_id uuid NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_feeds_user_id_user_id FOREIGN key (user_id) REFERENCES users (id) ON DELETE cascade,
    CONSTRAINT unique_feeds_name UNIQUE (name),
    CONSTRAINT unique_feeds_url UNIQUE (url)
);

-- +goose Down
DROP TABLE feeds;

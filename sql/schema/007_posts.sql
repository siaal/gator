-- +goose Up
CREATE TABLE posts (
    id uuid PRIMARY key,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title text NOT NULL,
    url text UNIQUE NOT NULL,
    description text NOT NULL,
    published_at TIMESTAMP NOT NULL,
    feed_id uuid NOT NULL,
    CONSTRAINT fk_posts_feed_id FOREIGN key (feed_id) REFERENCES feeds (id) ON DELETE cascade
);

-- +goose Down
DROP TABLE posts;

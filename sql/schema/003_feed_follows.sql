-- +goose Up
CREATE TABLE feed_follows (
    id uuid PRIMARY key,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id uuid,
    feed_id uuid,
    CONSTRAINT fk_feed_follows_user_id_users_id FOREIGN key (user_id) REFERENCES users (id),
    CONSTRAINT fk_feed_follows_feed_id_feeds_id FOREIGN key (feed_id) REFERENCES feeds (id),
    CONSTRAINT unq_feed_follows_user_id_feed_id UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;

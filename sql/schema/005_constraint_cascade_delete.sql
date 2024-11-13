-- +goose Up
ALTER TABLE feed_follows
DROP CONSTRAINT fk_feed_follows_user_id_users_id;

ALTER TABLE feed_follows add CONSTRAINT fk_feed_follows_user_id_users_id FOREIGN key (user_id) REFERENCES users (id) ON DELETE cascade;

ALTER TABLE feed_follows
DROP CONSTRAINT fk_feed_follows_feed_id_feeds_id;

ALTER TABLE feed_follows add CONSTRAINT fk_feed_follows_feed_id_feeds_id FOREIGN key (feed_id) REFERENCES feeds (id) ON DELETE cascade;

-- +goose Down
ALTER TABLE feed_follows
DROP CONSTRAINT fk_feed_follows_user_id_users_id;

ALTER TABLE feed_follows
DROP CONSTRAINT fk_feed_follows_feed_id_feeds_id;

ALTER TABLE feed_follows add CONSTRAINT fk_feed_follows_user_id_users_id FOREIGN key (user_id) REFERENCES users (id);

ALTER TABLE feed_follows add CONSTRAINT fk_feed_follows_feed_id_feeds_id FOREIGN key (feed_id) REFERENCES feeds (id);

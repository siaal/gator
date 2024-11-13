-- +goose Up
CREATE index idx_feeds_url ON feeds (url);

CREATE index idx_users_name ON users (name);

-- +goose Down
DROP index idx_feeds_url;

DROP index idx_users_name;

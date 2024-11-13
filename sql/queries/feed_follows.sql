-- name: CreateFeedFollow :one
WITH insert_ff AS (
    INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $3,
        $2,
        (SELECT id FROM users where users.name = @username), 
        (SELECT id FROM feeds WHERE feeds.url = @feed_url)
    )
    RETURNING *
)
SELECT ff.id, ff.created_at, ff.updated_at, ff.user_id, ff.feed_id, users.name, feeds.name
FROM insert_ff ff
    INNER JOIN users ON users.id = ff.user_id
    INNER JOIN feeds ON feeds.id = ff.feed_id;

-- name: GetFollowing :many
SELECT feed.id, feed.name, feed.url
FROM feeds feed
    INNER JOIN feed_follows ff ON ff.feed_id = feed.id
    INNER JOIN users ON users.id = ff.user_id
WHERE users.name = @username;

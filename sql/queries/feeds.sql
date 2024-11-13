-- name: AddFeed :one
INSERT INTO feeds(id, name, user_id, url, created_at, updated_at)
VALUES($1, $2, (SELECT id from users WHERE users.name = @username), $3, $4, $5)
RETURNING *;

-- name: GetFeeds :many
SELECT
  feed.id AS id,
  feed.name,
  feed.url,
  feed.created_at,
  feed.updated_at,
  users.id AS user_id,
  users.name AS user_name
FROM
  feeds feed
  INNER JOIN users 
  ON feed.user_id = users.id;

-- name: ClearFeeds :exec
DELETE FROM feeds;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = @now, updated_at = @now
WHERE id = @id;


-- name: NextToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST;

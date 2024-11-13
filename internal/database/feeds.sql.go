// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: feeds.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const addFeed = `-- name: AddFeed :one
INSERT INTO feeds(id, name, user_id, url, created_at, updated_at)
VALUES($1, $2, (SELECT id from users WHERE users.name = $6), $3, $4, $5)
RETURNING id, name, url, user_id, created_at, updated_at
`

type AddFeedParams struct {
	ID        uuid.UUID
	Name      string
	Url       string
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
}

func (q *Queries) AddFeed(ctx context.Context, arg AddFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, addFeed,
		arg.ID,
		arg.Name,
		arg.Url,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Username,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const clearFeeds = `-- name: ClearFeeds :exec
DELETE FROM feeds
`

func (q *Queries) ClearFeeds(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, clearFeeds)
	return err
}

const getFeeds = `-- name: GetFeeds :many
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
  ON feed.user_id = users.id
`

type GetFeedsRow struct {
	ID        uuid.UUID
	Name      string
	Url       string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	UserName  string
}

func (q *Queries) GetFeeds(ctx context.Context) ([]GetFeedsRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedsRow
	for rows.Next() {
		var i GetFeedsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Url,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.UserName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

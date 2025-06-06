// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: my_list.sql

package sqlcgen

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const deleteMyList = `-- name: DeleteMyList :exec
DELETE
FROM my_lists
WHERE id = $1
`

func (q *Queries) DeleteMyList(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteMyList, id)
	return err
}

const existsMyList = `-- name: ExistsMyList :one
SELECT EXISTS (
  SELECT 1 FROM my_lists WHERE id = $1
) AS exists
`

func (q *Queries) ExistsMyList(ctx context.Context, id int32) (bool, error) {
	row := q.db.QueryRowContext(ctx, existsMyList, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const getMyListByID = `-- name: GetMyListByID :one
SELECT id, user_id, name, position, created_at, updated_at FROM my_lists WHERE id = $1
`

func (q *Queries) GetMyListByID(ctx context.Context, id int32) (MyList, error) {
	row := q.db.QueryRowContext(ctx, getMyListByID, id)
	var i MyList
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Position,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const insertMyList = `-- name: InsertMyList :one
INSERT INTO my_lists (user_id, name, position, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, name, position, created_at, updated_at
`

type InsertMyListParams struct {
	UserID    uuid.NullUUID
	Name      string
	Position  int32
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}

func (q *Queries) InsertMyList(ctx context.Context, arg InsertMyListParams) (MyList, error) {
	row := q.db.QueryRowContext(ctx, insertMyList,
		arg.UserID,
		arg.Name,
		arg.Position,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i MyList
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Position,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listMyListsByUserID = `-- name: ListMyListsByUserID :many
SELECT id, user_id, name, position, created_at, updated_at FROM my_lists WHERE user_id = $1 ORDER BY id
`

func (q *Queries) ListMyListsByUserID(ctx context.Context, userID uuid.NullUUID) ([]MyList, error) {
	rows, err := q.db.QueryContext(ctx, listMyListsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MyList
	for rows.Next() {
		var i MyList
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.Position,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const updateMyListName = `-- name: UpdateMyListName :exec
UPDATE my_lists
SET name = $1,
    updated_at = $2
WHERE id = $3
`

type UpdateMyListNameParams struct {
	Name      string
	UpdatedAt sql.NullTime
	ID        int32
}

func (q *Queries) UpdateMyListName(ctx context.Context, arg UpdateMyListNameParams) error {
	_, err := q.db.ExecContext(ctx, updateMyListName, arg.Name, arg.UpdatedAt, arg.ID)
	return err
}

const updateMyListPosition = `-- name: UpdateMyListPosition :exec
UPDATE my_lists
SET position = $1,
    updated_at = $2
WHERE id = $3
`

type UpdateMyListPositionParams struct {
	Position  int32
	UpdatedAt sql.NullTime
	ID        int32
}

func (q *Queries) UpdateMyListPosition(ctx context.Context, arg UpdateMyListPositionParams) error {
	_, err := q.db.ExecContext(ctx, updateMyListPosition, arg.Position, arg.UpdatedAt, arg.ID)
	return err
}

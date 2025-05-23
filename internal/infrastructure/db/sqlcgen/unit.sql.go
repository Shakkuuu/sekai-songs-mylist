// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: unit.sql

package sqlcgen

import (
	"context"
)

const existsUnit = `-- name: ExistsUnit :one
SELECT EXISTS (
  SELECT 1 FROM units WHERE id = $1
) AS exists
`

func (q *Queries) ExistsUnit(ctx context.Context, id int32) (bool, error) {
	row := q.db.QueryRowContext(ctx, existsUnit, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const getUnitByID = `-- name: GetUnitByID :one
SELECT id, name FROM units WHERE id = $1
`

func (q *Queries) GetUnitByID(ctx context.Context, id int32) (Unit, error) {
	row := q.db.QueryRowContext(ctx, getUnitByID, id)
	var i Unit
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const insertUnit = `-- name: InsertUnit :one
INSERT INTO units (name)
VALUES ($1)
RETURNING id, name
`

func (q *Queries) InsertUnit(ctx context.Context, name string) (Unit, error) {
	row := q.db.QueryRowContext(ctx, insertUnit, name)
	var i Unit
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const listUnits = `-- name: ListUnits :many
SELECT id, name FROM units ORDER BY id
`

func (q *Queries) ListUnits(ctx context.Context) ([]Unit, error) {
	rows, err := q.db.QueryContext(ctx, listUnits)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Unit
	for rows.Next() {
		var i Unit
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
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

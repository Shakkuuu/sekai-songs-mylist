-- name: GetUnitByID :one
SELECT * FROM units WHERE id = $1;

-- name: ListUnits :many
SELECT * FROM units ORDER BY id;

-- name: InsertUnit :one
INSERT INTO units (name)
VALUES ($1)
RETURNING *;

-- name: ExistsUnit :one
SELECT EXISTS (
  SELECT 1 FROM units WHERE id = $1
) AS exists;

-- name: GetUnitByID :one
SELECT * FROM units WHERE id = $1;

-- name: ListUnits :many
SELECT * FROM units ORDER BY id;

-- name: InsertUnit :exec
INSERT INTO units (name)
VALUES ($1);

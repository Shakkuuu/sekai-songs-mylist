-- name: GetSingerByID :one
SELECT * FROM singers WHERE id = $1;

-- name: ListSingers :many
SELECT * FROM singers ORDER BY id;

-- name: InsertSinger :one
INSERT INTO singers (name)
VALUES ($1)
RETURNING *;

-- name: ExistsSinger :one
SELECT EXISTS (
  SELECT 1 FROM singers WHERE id = $1
) AS exists;

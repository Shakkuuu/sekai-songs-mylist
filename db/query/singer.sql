-- name: GetSingerByID :one
SELECT * FROM singers WHERE id = $1;

-- name: ListSingers :many
SELECT * FROM singers ORDER BY id;

-- name: InsertSinger :exec
INSERT INTO singers (name)
VALUES ($1);

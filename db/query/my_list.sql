-- name: GetMyListByID :one
SELECT * FROM my_lists WHERE id = $1;

-- name: ListMyListsByUserID :many
SELECT * FROM my_lists WHERE user_id = $1 ORDER BY id;

-- name: InsertMyList :one
INSERT INTO my_lists (user_id, name, position, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ExistsMyList :one
SELECT EXISTS (
  SELECT 1 FROM my_lists WHERE id = $1
) AS exists;

-- name: UpdateMyListName :exec
UPDATE my_lists
SET name = $1,
    updated_at = $2
WHERE id = $3;

-- name: UpdateMyListPosition :exec
UPDATE my_lists
SET position = $1,
    updated_at = $2
WHERE id = $3;

-- name: DeleteMyList :exec
DELETE
FROM my_lists
WHERE id = $1;

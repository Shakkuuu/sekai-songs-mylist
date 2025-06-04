-- name: GetMyListChartByID :one
SELECT * FROM my_list_charts WHERE id = $1;

-- name: ListMyListChartsByMyListID :many
SELECT * FROM my_list_charts WHERE my_list_id = $1 ORDER BY id;

-- name: InsertMyListChart :one
INSERT INTO my_list_charts (my_list_id, chart_id, clear_type, memo, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ExistsMyListChart :one
SELECT EXISTS (
    SELECT 1 FROM my_list_charts WHERE id = $1
) AS exists;

-- name: ExistsMyListChartByMyListIDAndChartID :one
SELECT EXISTS (
    SELECT 1 FROM my_list_charts WHERE my_list_id = $1 AND chart_id = $2
) AS exists;

-- name: UpdateMyListChartClearType :exec
UPDATE my_list_charts
SET clear_type = $1,
    updated_at = $2
WHERE id = $3;

-- name: UpdateMyListChartMemo :exec
UPDATE my_list_charts
SET memo = $1,
    updated_at = $2
WHERE id = $3;

-- name: DeleteMyListChart :exec
DELETE
FROM my_list_charts
WHERE id = $1;

-- name: DeleteMyListChartByMyListID :exec
DELETE
FROM my_list_charts
WHERE my_list_id = $1;

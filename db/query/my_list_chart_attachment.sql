-- name: GetMyListChartAttachmentByID :one
SELECT * FROM my_list_chart_attachments WHERE id = $1;

-- name: ListMyListChartAttachmentsByMyListChartID :many
SELECT * FROM my_list_chart_attachments WHERE my_list_chart_id = $1 ORDER BY id;

-- name: InsertMyListChartAttachment :one
INSERT INTO my_list_chart_attachments (my_list_chart_id, attachment_type, file_url, caption, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ExistsMyListChartAttachment :one
SELECT EXISTS (
    SELECT 1 FROM my_list_chart_attachments WHERE id = $1
) AS exists;

-- name: DeleteMyListChartAttachment :exec
DELETE
FROM my_list_chart_attachments
WHERE id = $1;

-- name: DeleteMyListChartAttachmentByMyListChartID :exec
DELETE
FROM my_list_chart_attachments
WHERE my_list_chart_id = $1;

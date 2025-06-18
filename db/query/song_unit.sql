-- name: InsertSongUnit :one
INSERT INTO song_units (song_id, unit_id)
VALUES ($1, $2)
RETURNING *;

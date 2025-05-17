-- name: GetArtistByID :one
SELECT * FROM artists WHERE id = $1;

-- name: ListArtists :many
SELECT * FROM artists ORDER BY id;

-- name: InsertArtist :one
INSERT INTO artists (name, kana)
VALUES ($1, $2)
RETURNING *;

-- name: ExistsArtist :one
SELECT EXISTS (
  SELECT 1 FROM artists WHERE id = $1
) AS exists;

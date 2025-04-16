-- name: GetArtistByID :one
SELECT * FROM artists WHERE id = $1;

-- name: ListArtists :many
SELECT * FROM artists ORDER BY id;

-- name: InsertArtist :exec
INSERT INTO artists (name, kana)
VALUES ($1, $2);

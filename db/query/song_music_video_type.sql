-- name: InsertSongMusicVideoType :one
INSERT INTO song_music_video_types (song_id, music_video_type)
VALUES ($1, $2)
RETURNING *;

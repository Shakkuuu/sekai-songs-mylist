-- name: InsertVocalPatterns :exec
INSERT INTO vocal_patterns (song_id, name)
VALUES ($1, $2);

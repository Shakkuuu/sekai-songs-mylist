-- name: InsertVocalPatternUnit :one
INSERT INTO vocal_pattern_units (vocal_pattern_id, unit_id)
VALUES ($1, $2)
RETURNING *;

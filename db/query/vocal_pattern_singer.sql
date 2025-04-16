-- name: InsertVocalPatternSingers :exec
INSERT INTO vocal_pattern_singers (vocal_pattern_id, singer_id, position)
VALUES ($1, $2, $3);

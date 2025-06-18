-- name: GetSongDetailsByID :many
SELECT
    s.id AS id,
    s.name,
    s.kana,

    l.id AS lyrics_artist_id,
    l.name AS lyrics_artist_name,
    l.kana AS lyrics_artist_kana,

    m.id AS music_artist_id,
    m.name AS music_artist_name,
    m.kana AS music_artist_kana,

    a.id AS arrangement_artist_id,
    a.name AS arrangement_artist_name,
    a.kana AS arrangement_artist_kana,

    s.thumbnail,
    s.original_video,
    s.release_time,
    s.deleted,

    vp.id AS vocal_pattern_id,
    vp.name AS vocal_pattern_name,

    vps.singer_id,
    si.name AS singer_name,
    vps.position AS singer_order,

    su.unit_id,
    u.name AS unit_name,

    smvt.music_video_type AS music_video_type_id

FROM songs s
LEFT JOIN artists l ON s.lyrics_id = l.id
LEFT JOIN artists m ON s.music_id = m.id
LEFT JOIN artists a ON s.arrangement_id = a.id
LEFT JOIN vocal_patterns vp ON vp.song_id = s.id
LEFT JOIN vocal_pattern_singers vps ON vps.vocal_pattern_id = vp.id
LEFT JOIN singers si ON vps.singer_id = si.id
LEFT JOIN song_units su ON su.song_id = s.id
LEFT JOIN units u ON su.unit_id = u.id
LEFT JOIN song_music_video_types smvt ON smvt.song_id = s.id
WHERE s.id = $1;

-- name: ListSongWithArtists :many
SELECT
    s.id AS id,
    s.name,
    s.kana,

    l.id AS lyrics_artist_id,
    l.name AS lyrics_artist_name,
    l.kana AS lyrics_artist_kana,

    m.id AS music_artist_id,
    m.name AS music_artist_name,
    m.kana AS music_artist_kana,

    a.id AS arrangement_artist_id,
    a.name AS arrangement_artist_name,
    a.kana AS arrangement_artist_kana,

    s.thumbnail,
    s.original_video,
    s.release_time,
    s.deleted,

    vp.id AS vocal_pattern_id,
    vp.name AS vocal_pattern_name,

    vps.singer_id,
    si.name AS singer_name,
    vps.position AS singer_order,

    su.unit_id,
    u.name AS unit_name,

    smvt.music_video_type AS music_video_type_id

FROM songs s
LEFT JOIN artists l ON s.lyrics_id = l.id
LEFT JOIN artists m ON s.music_id = m.id
LEFT JOIN artists a ON s.arrangement_id = a.id
LEFT JOIN vocal_patterns vp ON vp.song_id = s.id
LEFT JOIN vocal_pattern_singers vps ON vps.vocal_pattern_id = vp.id
LEFT JOIN singers si ON vps.singer_id = si.id
LEFT JOIN song_units su ON su.song_id = s.id
LEFT JOIN units u ON su.unit_id = u.id
LEFT JOIN song_music_video_types smvt ON smvt.song_id = s.id
ORDER BY s.id;

-- name: InsertSong :one
INSERT INTO songs (name, kana, lyrics_id, music_id, arrangement_id, thumbnail, original_video, release_time, deleted)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: ExistsSong :one
SELECT EXISTS (
  SELECT 1 FROM songs WHERE id = $1
) AS exists;

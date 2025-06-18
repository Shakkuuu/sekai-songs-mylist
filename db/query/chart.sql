-- name: GetChartWithSongWithArtistsByID :many
SELECT
    c.id AS id,

    s.id AS song_id,
    s.name AS song_name,
    s.kana AS song_kana,

    la.id AS lyrics_artist_id,
    la.name AS lyrics_artist_name,
    la.kana AS lyrics_artist_kana,

    ma.id AS music_artist_id,
    ma.name AS music_artist_name,
    ma.kana AS music_artist_kana,

    aa.id AS arrangement_artist_id,
    aa.name AS arrangement_artist_name,
    aa.kana AS arrangement_artist_kana,

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

    smvt.music_video_type AS music_video_type_id,

    c.difficulty_type,
    c.level,
    c.chart_view_link

FROM charts c
LEFT JOIN songs s ON c.song_id = s.id
LEFT JOIN artists la ON s.lyrics_id = la.id
LEFT JOIN artists ma ON s.music_id = ma.id
LEFT JOIN artists aa ON s.arrangement_id = aa.id
LEFT JOIN vocal_patterns vp ON vp.song_id = s.id
LEFT JOIN vocal_pattern_singers vps ON vps.vocal_pattern_id = vp.id
LEFT JOIN singers si ON vps.singer_id = si.id
LEFT JOIN song_units su ON su.song_id = s.id
LEFT JOIN units u ON su.unit_id = u.id
LEFT JOIN song_music_video_types smvt ON smvt.song_id = s.id
WHERE c.id = $1;

-- name: ListChartWithSongWithArtists :many
SELECT
    c.id AS id,

    s.id AS song_id,
    s.name AS song_name,
    s.kana AS song_kana,

    la.id AS lyrics_artist_id,
    la.name AS lyrics_artist_name,
    la.kana AS lyrics_artist_kana,

    ma.id AS music_artist_id,
    ma.name AS music_artist_name,
    ma.kana AS music_artist_kana,

    aa.id AS arrangement_artist_id,
    aa.name AS arrangement_artist_name,
    aa.kana AS arrangement_artist_kana,

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

    smvt.music_video_type AS music_video_type_id,

    c.difficulty_type,
    c.level,
    c.chart_view_link

FROM charts c
LEFT JOIN songs s ON c.song_id = s.id
LEFT JOIN artists la ON s.lyrics_id = la.id
LEFT JOIN artists ma ON s.music_id = ma.id
LEFT JOIN artists aa ON s.arrangement_id = aa.id
LEFT JOIN vocal_patterns vp ON vp.song_id = s.id
LEFT JOIN vocal_pattern_singers vps ON vps.vocal_pattern_id = vp.id
LEFT JOIN singers si ON vps.singer_id = si.id
LEFT JOIN song_units su ON su.song_id = s.id
LEFT JOIN units u ON su.unit_id = u.id
LEFT JOIN song_music_video_types smvt ON smvt.song_id = s.id
ORDER BY c.id;

-- name: InsertChart :one
INSERT INTO charts (song_id, difficulty_type, level, chart_view_link)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ExistsChart :one
SELECT EXISTS (
  SELECT 1 FROM charts WHERE id = $1
) AS exists;

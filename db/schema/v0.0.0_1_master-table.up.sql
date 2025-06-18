CREATE TABLE artists (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    kana VARCHAR(255) NOT NULL
);

CREATE TABLE singers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE units (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    kana VARCHAR(255) NOT NULL,
    lyrics_id INT REFERENCES artists(id),
    music_id INT REFERENCES artists(id),
    arrangement_id INT REFERENCES artists(id),
    thumbnail TEXT,
    original_video TEXT,
    release_time TIMESTAMP,
    deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE charts (
    id SERIAL PRIMARY KEY,
    song_id INT REFERENCES songs(id),
    difficulty_type INT,
    level INT,
    chart_view_link TEXT
);

CREATE TABLE vocal_patterns (
    id SERIAL PRIMARY KEY,
    song_id INT REFERENCES songs(id),
    name VARCHAR(255)
);

CREATE TABLE vocal_pattern_singers (
    id SERIAL PRIMARY KEY,
    vocal_pattern_id INT REFERENCES vocal_patterns(id),
    singer_id INT REFERENCES singers(id),
    position INT
);

CREATE TABLE song_units (
    id SERIAL PRIMARY KEY,
    song_id INT REFERENCES songs(id),
    unit_id INT REFERENCES units(id)
);

CREATE TABLE song_music_video_types (
    id SERIAL PRIMARY KEY,
    song_id INT REFERENCES songs(id),
    music_video_type INT
);

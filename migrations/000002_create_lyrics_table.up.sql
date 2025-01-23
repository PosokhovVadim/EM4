CREATE TABLE lyrics(
    song_id INTEGER NOT NULL,
    verse_number INTEGER NOT NULL,
    text TEXT NOT NULL,
    CHECK (verse_number >= 0),
    PRIMARY KEY (song_id, verse_number),
    FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE
);
package model

import "time"

type Lyrics struct {
	SongID      uint   `json:"song_id"`
	VerseNumber uint   `json:"verse_number"`
	Text        string `json:"text"`
}

type Song struct {
	ID          uint      `json:"id"`
	Group       string    `json:"group"`
	Name        string    `json:"name"`
	ReleaseDate time.Time `json:"release_date"`
	Link        string    `json:"link"`
	InsertedAt  time.Time `json:"inserted_at"`
}

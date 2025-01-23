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

type SongUpdate struct {
	Group       string          `json:"group,omitempty"`
	Name        string          `json:"name,omitempty"`
	ReleaseDate string          `json:"release_date,omitempty"`
	Link        string          `json:"link,omitempty"`
	Verses      map[uint]string `json:"verses,omitempty"`
}

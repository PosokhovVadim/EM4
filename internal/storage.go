package storage

import (
	"em4/internal/model"
	"errors"
)

var (
	ErrSongNotFound = errors.New("song not found")
)

// just for me
type Storage interface {
	AddSong(song model.Song, verses []string) (uint, error)
	DeleteSong(songID uint) error
	GetLyrics(songID uint, limit, offset int) ([]model.Lyrics, error)
	GetSong(songID uint) (*model.Song, error)
	GetAllSongs(filters map[string]string, limit, offset int) ([]model.Song, error)
	GetAllSongLyrics(songID uint) ([]model.Lyrics, error)
	UpdateSong(songID uint, updates model.SongUpdate) error
}

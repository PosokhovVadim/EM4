package storage

import (
	"context"
	"database/sql"
	storage "em4/internal"
	"em4/internal/model"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"

	"github.com/jackc/pgx/v4/pgxpool"
	"gorm.io/driver/postgres"
)

// first type
type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(path string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", path)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{
		db: db,
	}, nil
}

func (s *PostgresStorage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// second type
type ORMPostgresStorage struct {
	db *gorm.DB
}

func NewORMPostgresStorage(path string) (*ORMPostgresStorage, error) {
	db, err := gorm.Open(postgres.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&model.Song{}, &model.Lyrics{})

	return &ORMPostgresStorage{
		db: db,
	}, nil
}

func (s *ORMPostgresStorage) Close() error {
	if s.db != nil {
		sqlDB, err := s.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// third type
type PGXStorage struct {
	db *pgxpool.Pool
}

func NewPGXStorage(path string) (*PGXStorage, error) {
	config, _ := pgxpool.ParseConfig(path)
	config.MaxConns = 10

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal("Unable to create connection pool:", err)
	}

	return &PGXStorage{
		db: pool,
	}, nil
}

func (s *PGXStorage) Close() error {
	if s.db == nil {
		return nil
	}
	s.db.Close()
	return nil
}

// examples:
func (s *PGXStorage) GetSong(songID uint) (uint, error) {
	err := s.db.QueryRow(context.Background(), "SELECT id FROM songs WHERE id = $1", songID).Scan(&songID)
	return songID, err
}

func (s *ORMPostgresStorage) GetSong(songID uint) (*model.Song, error) {
	var song model.Song
	err := s.db.First(&song, songID).Error
	return &song, err
}

// main db
func (s *PostgresStorage) WithTransaction(fn func(tx *sql.Tx) error) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		err = tx.Rollback()
		if err != nil {
			return err
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) AddSong(song model.Song, verses []string) (uint, error) {
	var songID uint

	if err := s.WithTransaction(func(tx *sql.Tx) error {
		err := tx.QueryRow(
			`INSERT INTO songs (group_name, name, link, release_date, inserted_at) 
             VALUES ($1, $2, $3, $4, NOW()) 
             RETURNING id`,
			song.Group, song.Name, song.Link, song.ReleaseDate,
		).Scan(&songID)
		if err != nil {
			return err
		}

		for i, verse := range verses {
			_, err = tx.Exec(
				`INSERT INTO lyrics (song_id, verse_number, text) 
                 VALUES ($1, $2, $3)`,
				songID, i+1, verse,
			)
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return 0, err
	}

	return songID, nil
}

func (s *PostgresStorage) DeleteSong(songID uint) error {
	result, err := s.db.Exec(
		`DELETE FROM songs WHERE id = $1`,
		songID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return storage.ErrSongNotFound
	}

	return nil
}

func (s *PostgresStorage) GetSong(songID uint) (*model.Song, error) {
	song := &model.Song{}
	err := s.db.QueryRow(
		`SELECT id, group_name, name, link, release_date, inserted_at 
         FROM songs 
         WHERE id = $1`,
		songID,
	).Scan(
		&song.ID, &song.Group, &song.Name, &song.Link, &song.ReleaseDate, &song.InsertedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrSongNotFound
	}
	if err != nil {
		return nil, err
	}

	return song, nil
}

func (s *PostgresStorage) UpdateSong(songID uint, updates model.SongUpdate) error {
	songQuery, songArgs, err := s.buildUpdateSongQuery(songID, updates)
	if err != nil && len(updates.Verses) == 0 {
		return err
	}

	if songQuery != "" {
		_, err := s.db.Exec(songQuery, songArgs...)
		if err != nil {
			return fmt.Errorf("failed to update song: %w", err)
		}
	}

	verseQueries := s.buildUpdateVerseQuery(songID, updates.Verses)

	for _, q := range verseQueries {
		_, err := s.db.Exec(q.Query, q.Args...)
		if err != nil {
			return fmt.Errorf("failed to update verse: %w", err)
		}
	}
	return nil
}

func (s *PostgresStorage) buildUpdateSongQuery(songID uint, updates model.SongUpdate) (string, []interface{}, error) {
	query := "UPDATE songs SET "
	var args []interface{}
	argIndex := 1
	updatesApplied := false

	if updates.Group != "" {
		query += fmt.Sprintf("group_name = $%d, ", argIndex)
		args = append(args, updates.Group)
		argIndex++
		updatesApplied = true
	}
	if updates.Name != "" {
		query += fmt.Sprintf("name = $%d, ", argIndex)
		args = append(args, updates.Name)
		argIndex++
		updatesApplied = true
	}
	if updates.ReleaseDate != "" {
		query += fmt.Sprintf("release_date = $%d, ", argIndex)
		args = append(args, updates.ReleaseDate)
		argIndex++
		updatesApplied = true
	}
	if updates.Link != "" {
		query += fmt.Sprintf("link = $%d, ", argIndex)
		args = append(args, updates.Link)
		argIndex++
		updatesApplied = true
	}
	if !updatesApplied {
		return "", nil, fmt.Errorf("no valid fields to update")
	}

	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE id = $%d RETURNING id", argIndex)
	args = append(args, songID)
	return query, args, nil
}

func (s *PostgresStorage) buildUpdateVerseQuery(songID uint, verses map[uint]string) []struct {
	Query string
	Args  []interface{}
} {
	var queries []struct {
		Query string
		Args  []interface{}
	}

	for verseNumber, text := range verses {
		query := "UPDATE lyrics SET text = $1 WHERE song_id = $2 AND verse_number = $3"
		args := []interface{}{text, songID, verseNumber}
		queries = append(queries, struct {
			Query string
			Args  []interface{}
		}{
			Query: query,
			Args:  args,
		})
	}

	return queries
}

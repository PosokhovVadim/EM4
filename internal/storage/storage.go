package storage

import (
	"context"
	"database/sql"
	"em4/internal/model"
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


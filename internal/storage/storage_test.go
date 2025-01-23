package storage

import (
	"database/sql"
	st "em4/internal"
	"errors"
	"testing"
	"time"

	"em4/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddSong_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	storage := &PostgresStorage{
		db: db,
	}

	song := model.Song{
		Group:       "Test Group",
		Name:        "Test Song",
		Link:        "http://example.com",
		ReleaseDate: time.Now(),
	}
	verses := []string{"Verse 1", "Verse 2", "Verse 3"}

	mock.ExpectBegin()

	mock.ExpectQuery(`INSERT INTO songs \(group_name, name, link, release_date, inserted_at\) VALUES \(\$1, \$2, \$3, \$4, NOW\(\)\) RETURNING id`).
		WithArgs(song.Group, song.Name, song.Link, song.ReleaseDate).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	for i, verse := range verses {
		mock.ExpectExec(`INSERT INTO lyrics \(song_id, verse_number, text\) VALUES \(\$1, \$2, \$3\)`).
			WithArgs(1, i+1, verse).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	mock.ExpectCommit()

	songID, err := storage.AddSong(song, verses)
	require.NoError(t, err)

	assert.Equal(t, uint(1), songID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSong(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	storage := &PostgresStorage{db: db}
	songID := uint(1)

	expectedSong := &model.Song{
		ID:          songID,
		Group:       "Test Group",
		Name:        "Test Song",
		Link:        "http://example.com",
		ReleaseDate: time.Now(),
		InsertedAt:  time.Now(),
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, group_name, name, link, release_date, inserted_at FROM songs WHERE id = \$1`).
			WithArgs(songID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "group_name", "name", "link", "release_date", "inserted_at"}).
				AddRow(expectedSong.ID, expectedSong.Group, expectedSong.Name, expectedSong.Link, expectedSong.ReleaseDate, expectedSong.InsertedAt))

		song, err := storage.GetSong(songID)
		require.NoError(t, err)
		require.Equal(t, expectedSong, song)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, group_name, name, link, release_date, inserted_at FROM songs WHERE id = \$1`).
			WithArgs(songID).
			WillReturnError(sql.ErrNoRows)

		song, err := storage.GetSong(songID)
		require.ErrorIs(t, err, st.ErrSongNotFound)
		require.Nil(t, song)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, group_name, name, link, release_date, inserted_at FROM songs WHERE id = \$1`).
			WithArgs(songID).
			WillReturnError(errors.New("mocked database error"))

		song, err := storage.GetSong(songID)
		require.Error(t, err)
		require.Contains(t, err.Error(), "mocked database error")
		require.Nil(t, song)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteSong(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	storage := &PostgresStorage{db: db}
	songID := uint(1)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM songs WHERE id = \$1`).
			WithArgs(songID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := storage.DeleteSong(songID)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM songs WHERE id = \$1`).
			WithArgs(songID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := storage.DeleteSong(songID)
		require.ErrorIs(t, err, st.ErrSongNotFound)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM songs WHERE id = \$1`).
			WithArgs(songID).
			WillReturnError(errors.New("mocked database error"))

		err := storage.DeleteSong(songID)
		require.Error(t, err)
		require.Contains(t, err.Error(), "mocked database error")
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

package services

import (
	"testing"
	"time"

	"project-MVP/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTagStore_CreateTag_Success проверяет создание тега.
func TestTagStore_CreateTag_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTagStore(db)

	mock.ExpectQuery(`INSERT INTO tags \(name, created_at\) VALUES \(\$1,\$2\) RETURNING id`).
		WithArgs("urgent", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	tag, err := store.CreateTag(models.Tag{Name: "urgent"})

	require.NoError(t, err)
	assert.Equal(t, 1, tag.ID)
	assert.Equal(t, "urgent", tag.Name)
	assert.WithinDuration(t, time.Now(), tag.CreatedAt, time.Second)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagStore_CreateTag_EmptyName проверяет отказ при пустом названии тега.
func TestTagStore_CreateTag_EmptyName(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTagStore(db)

	_, err = store.CreateTag(models.Tag{Name: "   "})
	assert.EqualError(t, err, "tag name is required")
}

// TestTagStore_GetTags_Success проверяет получение списка тегов.
func TestTagStore_GetTags_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTagStore(db)
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "name", "created_at"}).
		AddRow(1, "bug", createdAt).
		AddRow(2, "feature", createdAt)

	mock.ExpectQuery(`SELECT id, name, created_at FROM tags`).
		WillReturnRows(rows)

	tags, err := store.GetTags()

	require.NoError(t, err)
	assert.Len(t, tags, 2)
	assert.Equal(t, "bug", tags[0].Name)
	assert.Equal(t, "feature", tags[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagStore_GetTags_Empty проверяет обработку пустого списка тегов.
func TestTagStore_GetTags_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTagStore(db)

	mock.ExpectQuery(`SELECT id, name, created_at FROM tags`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at"}))

	tags, err := store.GetTags()

	require.NoError(t, err)
	assert.Empty(t, tags)
	assert.NoError(t, mock.ExpectationsWereMet())
}

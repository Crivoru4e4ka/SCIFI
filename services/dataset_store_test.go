package services

import (
	"database/sql"
	"testing"
	"time"

	"project-MVP/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDatasetStore_GetAllDatasets_Success проверяет получение всех датасетов.
func TestDatasetStore_GetAllDatasets_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewDatasetStore(db)
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "project_id", "name", "description", "version", "data_url", "parameters", "created_at"}).
		AddRow(1, 1, "DS1", "Desc", "v1", "http://example.com", nil, createdAt).
		AddRow(2, 1, "DS2", "Desc2", "v2", "http://example2.com", []byte(`{"key":"val"}`), createdAt)

	mock.ExpectQuery(`SELECT id, project_id, name, description, version, data_url, parameters, created_at FROM datasets`).
		WillReturnRows(rows)

	datasets, err := store.GetAllDatasets()

	require.NoError(t, err)
	assert.Len(t, datasets, 2)
	assert.Equal(t, "DS1", datasets[0].Name)
	assert.Equal(t, "v2", datasets[1].Version)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestDatasetStore_CreateDataset_Success проверяет создание датасета.
func TestDatasetStore_CreateDataset_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewDatasetStore(db)
	createdAt := time.Now()

	mock.ExpectQuery(`INSERT INTO datasets \(project_id, name, description, version, data_url, parameters\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\) RETURNING id, created_at`).
		WithArgs(1, "DS", "Desc", "v1", "http://example.com", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(5, sql.NullTime{Valid: true, Time: createdAt}))

	ds, err := store.CreateDataset(models.Dataset{
		ProjectID:   1,
		Name:        "DS",
		Description: "Desc",
		Version:     "v1",
		DataURL:     "http://example.com",
	})

	require.NoError(t, err)
	assert.Equal(t, 5, ds.ID)
	assert.WithinDuration(t, createdAt, ds.CreatedAt, time.Second)
	assert.NoError(t, mock.ExpectationsWereMet())
}

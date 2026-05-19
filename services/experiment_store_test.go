package services

import (
	"testing"

	"project-MVP/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestExperimentStore_GetAllExperimentDatasets_Success проверяет получение
// всех связей задач с датасетами.
func TestExperimentStore_GetAllExperimentDatasets_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewExperimentStore(db)

	rows := sqlmock.NewRows([]string{"task_id", "dataset_id"}).
		AddRow(1, 1).
		AddRow(1, 2).
		AddRow(2, 1)

	mock.ExpectQuery(`SELECT task_id, dataset_id FROM experiment_datasets`).
		WillReturnRows(rows)

	result, err := store.GetAllExperimentDatasets()

	require.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, 1, result[0].TaskID)
	assert.Equal(t, 2, result[1].DatasetID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestExperimentStore_CreateExperimentDataset_Success проверяет создание
// связи задачи и датасета.
func TestExperimentStore_CreateExperimentDataset_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewExperimentStore(db)

	mock.ExpectExec(`INSERT INTO experiment_datasets \(task_id, dataset_id\) VALUES \(\$1, \$2\)`).
		WithArgs(5, 3).
		WillReturnResult(sqlmock.NewResult(0, 1))

	eds, err := store.CreateExperimentDataset(models.ExperimentDataset{
		TaskID:    5,
		DatasetID: 3,
	})

	require.NoError(t, err)
	assert.Equal(t, 5, eds.TaskID)
	assert.Equal(t, 3, eds.DatasetID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

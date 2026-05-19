package services

import (
	"database/sql"
	"testing"

	"project-MVP/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSprintStore_GetProjectSprints_Success проверяет получение спринтов проекта.
func TestSprintStore_GetProjectSprints_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewSprintStore(db)

	rows := sqlmock.NewRows([]string{"id", "project_id", "name", "status", "start_date", "end_date"}).
		AddRow(1, 1, "Sprint 1", "active", nil, nil).
		AddRow(2, 1, "Sprint 2", "closed", nil, nil)

	mock.ExpectQuery(`SELECT id, project_id, name, status, start_date, end_date FROM sprints WHERE project_id = \$1 ORDER BY created_at DESC`).
		WithArgs(1).
		WillReturnRows(rows)

	sprints, err := store.GetProjectSprints(1)

	require.NoError(t, err)
	assert.Len(t, sprints, 2)
	assert.Equal(t, "Sprint 1", sprints[0].Name)
	assert.Equal(t, "active", sprints[0].Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSprintStore_GetProjectSprints_Empty проверяет, что пустой результат
// возвращается как пустой слайс.
func TestSprintStore_GetProjectSprints_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewSprintStore(db)

	mock.ExpectQuery(`SELECT id, project_id, name, status, start_date, end_date FROM sprints WHERE project_id = \$1 ORDER BY created_at DESC`).
		WithArgs(99).
		WillReturnRows(sqlmock.NewRows([]string{"id", "project_id", "name", "status", "start_date", "end_date"}))

	sprints, err := store.GetProjectSprints(99)

	require.NoError(t, err)
	assert.Empty(t, sprints)
	assert.NotNil(t, sprints)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSprintStore_CreateSprint_Success проверяет создание спринта.
func TestSprintStore_CreateSprint_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewSprintStore(db)

	mock.ExpectQuery(`INSERT INTO sprints \(project_id, name, status, start_date, end_date\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING id`).
		WithArgs(1, "New Sprint", "planned", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(5))

	sprint := models.Sprint{
		ProjectID: 1,
		Name:      "New Sprint",
		Status:    "planned",
	}
	err = store.CreateSprint(&sprint)

	require.NoError(t, err)
	assert.Equal(t, 5, sprint.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSprintStore_GetSprintByID_Success проверяет получение спринта по ID.
func TestSprintStore_GetSprintByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewSprintStore(db)

	rows := sqlmock.NewRows([]string{"id", "project_id", "name", "status", "start_date", "end_date", "goal"}).
		AddRow(1, 1, "S1", "active", nil, nil, "Goal")

	mock.ExpectQuery(`SELECT id, project_id, name, status, start_date, end_date, goal FROM sprints WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	sprint, err := store.GetSprintByID(1)

	require.NoError(t, err)
	assert.Equal(t, 1, sprint.ID)
	assert.Equal(t, "Goal", sprint.Goal)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSprintStore_GetSprintByID_NotFound проверяет обработку отсутствующего спринта.
func TestSprintStore_GetSprintByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewSprintStore(db)

	mock.ExpectQuery(`SELECT id, project_id, name, status, start_date, end_date, goal FROM sprints WHERE id = \$1`).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	_, err = store.GetSprintByID(999)

	assert.ErrorIs(t, err, ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSprintStore_StartSprint_Success проверяет активацию спринта.
func TestSprintStore_StartSprint_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewSprintStore(db)

	mock.ExpectExec(`UPDATE sprints SET name=\$1, start_date=\$2, end_date=\$3, goal=\$4, status='active' WHERE id=\$5`).
		WithArgs("Sprint X", sqlmock.AnyArg(), sqlmock.AnyArg(), "New Goal", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = store.StartSprint(1, models.Sprint{
		Name: "Sprint X",
		Goal: "New Goal",
	})

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSprintStore_CompleteSprint_Success проверяет завершение спринта
// с возвратом незавершенных задач в бэклог.
func TestSprintStore_CompleteSprint_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewSprintStore(db)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE sprints SET status = 'closed', end_date = NOW\(\) WHERE id = \$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`UPDATE tasks SET sprint_id = NULL WHERE sprint_id = \$1 AND status != 'done'`).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 3))
	mock.ExpectCommit()

	err = store.CompleteSprint(1)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

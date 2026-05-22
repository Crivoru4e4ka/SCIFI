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

// TestProjectStore_GetProjectByID_Success проверяет получение проекта по ID.
func TestProjectStore_GetProjectByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewProjectStore(db)
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{
		"id", "name", "key", "description", "research_goal", "main_hypothesis",
		"novelty", "expected_result", "visibility", "start_date", "end_date",
		"status", "created_by", "created_at", "team_id", "execution_type",
	}).AddRow(1, "AI Research", "AIR", "Desc", "Goal", "Hypothesis",
		"Novel", "Result", "open", startDate, nil,
		"active", 1, startDate, nil, "manual")

	mock.ExpectQuery(`SELECT id, name, key, description, research_goal, main_hypothesis, novelty, expected_result, visibility, start_date, end_date, status, created_by, created_at, team_id, execution_type FROM projects WHERE id=\$1`).
		WithArgs(1).
		WillReturnRows(rows)

	project, err := store.GetProjectByID(1)

	require.NoError(t, err)
	assert.Equal(t, 1, project.Id)
	assert.Equal(t, "AI Research", project.Name)
	assert.Equal(t, "AIR", project.Key)
	assert.Equal(t, "active", project.Status)
	assert.Nil(t, project.TeamId)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestProjectStore_GetProjectByID_NotFound проверяет обработку отсутствующего проекта.
func TestProjectStore_GetProjectByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewProjectStore(db)

	mock.ExpectQuery(`SELECT id, name, key, description, research_goal, main_hypothesis, novelty, expected_result, visibility, start_date, end_date, status, created_by, created_at, team_id, execution_type FROM projects WHERE id=\$1`).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	_, err = store.GetProjectByID(999)

	assert.ErrorIs(t, err, ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestProjectStore_GetProjectByID_WithTeamID проверяет корректное
// сканирование nullable team_id.
func TestProjectStore_GetProjectByID_WithTeamID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewProjectStore(db)
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{
		"id", "name", "key", "description", "research_goal", "main_hypothesis",
		"novelty", "expected_result", "visibility", "start_date", "end_date",
		"status", "created_by", "created_at", "team_id", "execution_type",
	}).AddRow(2, "Team Project", "TP", "Desc", "Goal", "Hyp",
		"Nov", "Res", "closed", startDate, nil,
		"active", 1, startDate, 5, "team")

	mock.ExpectQuery(`SELECT id, name, key, description, research_goal, main_hypothesis, novelty, expected_result, visibility, start_date, end_date, status, created_by, created_at, team_id, execution_type FROM projects WHERE id=\$1`).
		WithArgs(2).
		WillReturnRows(rows)

	project, err := store.GetProjectByID(2)

	require.NoError(t, err)
	require.NotNil(t, project.TeamId)
	assert.Equal(t, 5, *project.TeamId)
	assert.Equal(t, "team", project.ExecutionType)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestProjectStore_GetProjects_Success проверяет получение списка всех проектов.
func TestProjectStore_GetProjects_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewProjectStore(db)
	startDate := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "name", "key", "description", "research_goal", "main_hypothesis",
		"novelty", "expected_result", "visibility", "start_date", "end_date",
		"status", "created_by", "created_at", "team_id", "execution_type",
	}).
		AddRow(1, "P1", "K1", "D1", "G1", "H1", "N1", "E1", "open", startDate, nil, "active", 1, startDate, nil, "manual").
		AddRow(2, "P2", "K2", "D2", "G2", "H2", "N2", "E2", "closed", startDate, nil, "active", 2, startDate, nil, "manual")

	mock.ExpectQuery(`SELECT id, name, key, description, research_goal, main_hypothesis, novelty, expected_result, visibility, start_date, end_date, status, created_by, created_at, team_id, execution_type FROM projects`).
		WillReturnRows(rows)

	projects, err := store.GetProjects()

	require.NoError(t, err)
	assert.Len(t, projects, 2)
	assert.Equal(t, "P1", projects[0].Name)
	assert.Equal(t, "P2", projects[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestProjectStore_GetUserProjects_Success проверяет получение проектов
// пользователя с учетом visibility и участия.
func TestProjectStore_GetUserProjects_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewProjectStore(db)
	startDate := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "name", "key", "description", "research_goal", "main_hypothesis",
		"novelty", "expected_result", "visibility", "start_date", "end_date",
		"status", "created_by", "created_at", "team_id", "execution_type",
	}).AddRow(1, "My Project", "MP", "Desc", "Goal", "Hyp",
		"Nov", "Res", "closed", startDate, nil,
		"active", 7, startDate, nil, "manual")

	mock.ExpectQuery(`SELECT DISTINCT p.id, p.name, p.key, p.description, p.research_goal, p.main_hypothesis, p.novelty, p.expected_result, p.visibility, p.start_date, p.end_date, p.status, p.created_by, p.created_at, p.team_id, p.execution_type FROM projects p LEFT JOIN project_members pm ON p.id = pm.project_id WHERE p.created_by = \$1 OR pm.user_id = \$1 OR p.visibility = 'open' ORDER BY p.created_at DESC`).
		WithArgs(7).
		WillReturnRows(rows)

	projects, err := store.GetUserProjects(7)

	require.NoError(t, err)
	assert.Len(t, projects, 1)
	assert.Equal(t, "My Project", projects[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestProjectStore_GetUserProjects_EmptyResult проверяет, что пустой
// результат возвращается как пустой слайс (не nil).
func TestProjectStore_GetUserProjects_EmptyResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewProjectStore(db)

	mock.ExpectQuery(`SELECT DISTINCT p.id, p.name, p.key, p.description, p.research_goal, p.main_hypothesis, p.novelty, p.expected_result, p.visibility, p.start_date, p.end_date, p.status, p.created_by, p.created_at, p.team_id, p.execution_type FROM projects p LEFT JOIN project_members pm ON p.id = pm.project_id WHERE p.created_by = \$1 OR pm.user_id = \$1 OR p.visibility = 'open' ORDER BY p.created_at DESC`).
		WithArgs(99).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "key", "description", "research_goal", "main_hypothesis",
			"novelty", "expected_result", "visibility", "start_date", "end_date",
			"status", "created_by", "created_at", "team_id", "execution_type",
		}))

	projects, err := store.GetUserProjects(99)

	require.NoError(t, err)
	assert.Empty(t, projects)
	assert.NotNil(t, projects)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestProjectStore_GetProjectProgress_ZeroTasks проверяет, что
// прогресс равен 0, если в проекте нет задач.
func TestProjectStore_GetProjectProgress_ZeroTasks(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewProjectStore(db)
	startDate := time.Now()

	// GetProjectByID
	projRows := sqlmock.NewRows([]string{
		"id", "name", "key", "description", "research_goal", "main_hypothesis",
		"novelty", "expected_result", "visibility", "start_date", "end_date",
		"status", "created_by", "created_at", "team_id", "execution_type",
	}).AddRow(1, "P", "K", "D", "G", "H", "N", "E", "open", startDate, nil, "active", 1, startDate, nil, "manual")
	mock.ExpectQuery(`SELECT id, name, key, description, research_goal, main_hypothesis, novelty, expected_result, visibility, start_date, end_date, status, created_by, created_at, team_id, execution_type FROM projects WHERE id=\$1`).
		WithArgs(1).
		WillReturnRows(projRows)

	// Progress query
	mock.ExpectQuery(`SELECT COUNT\(CASE WHEN status = 'done' THEN 1 END\) AS done_count, COUNT\(\*\) AS total_count FROM tasks WHERE project_id=\$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"done_count", "total_count"}).AddRow(0, 0))

	progress, err := store.GetProjectProgress(1)

	require.NoError(t, err)
	assert.Equal(t, 0.0, progress)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestProjectStore_GetProjectProgress_WithTasks проверяет корректный
// расчет процента выполненных задач.
func TestProjectStore_GetProjectProgress_WithTasks(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewProjectStore(db)
	startDate := time.Now()

	// GetProjectByID
	projRows := sqlmock.NewRows([]string{
		"id", "name", "key", "description", "research_goal", "main_hypothesis",
		"novelty", "expected_result", "visibility", "start_date", "end_date",
		"status", "created_by", "created_at", "team_id", "execution_type",
	}).AddRow(1, "P", "K", "D", "G", "H", "N", "E", "open", startDate, nil, "active", 1, startDate, nil, "manual")
	mock.ExpectQuery(`SELECT id, name, key, description, research_goal, main_hypothesis, novelty, expected_result, visibility, start_date, end_date, status, created_by, created_at, team_id, execution_type FROM projects WHERE id=\$1`).
		WithArgs(1).
		WillReturnRows(projRows)

	// Progress query: 3 done out of 12
	mock.ExpectQuery(`SELECT COUNT\(CASE WHEN status = 'done' THEN 1 END\) AS done_count, COUNT\(\*\) AS total_count FROM tasks WHERE project_id=\$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"done_count", "total_count"}).AddRow(3, 12))

	progress, err := store.GetProjectProgress(1)

	require.NoError(t, err)
	assert.InDelta(t, 25.0, progress, 0.01)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestProjectStore_GetProjectProgress_ProjectNotFound проверяет,
// что при отсутствии проекта возвращается ErrNotFound.
func TestProjectStore_GetProjectProgress_ProjectNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewProjectStore(db)

	mock.ExpectQuery(`SELECT id, name, key, description, research_goal, main_hypothesis, novelty, expected_result, visibility, start_date, end_date, status, created_by, created_at, team_id, execution_type FROM projects WHERE id=\$1`).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	_, err = store.GetProjectProgress(999)

	assert.ErrorIs(t, err, ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestProjectStore_UpdateProject_Success проверяет обновление проекта.
func TestProjectStore_UpdateProject_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewProjectStore(db)

	mock.ExpectExec(`UPDATE projects SET name = \$1, description = \$2, research_goal = \$3, main_hypothesis = \$4, novelty = \$5, expected_result = \$6, status = \$7, end_date = \$8, visibility = \$9 WHERE id = \$10`).
		WithArgs("New Name", "New Desc", "Goal", "Hyp", "Nov", "Res", "active", sqlmock.AnyArg(), "open", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = store.UpdateProject(1, models.Project{
		Name:        "New Name",
		Description: "New Desc",
		ResearchGoal: "Goal",
		MainHypothesis: "Hyp",
		Novelty: "Nov",
		ExpectedResult: "Res",
		Status: "active",
		Visibility: "open",
	})

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestProjectStore_DeleteProject_Success проверяет каскадное удаление
// проекта через транзакцию.
func TestProjectStore_DeleteProject_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewProjectStore(db)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM tasks WHERE project_id = \$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 5))
	mock.ExpectExec(`DELETE FROM project_members WHERE project_id = \$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 3))
	mock.ExpectExec(`DELETE FROM projects WHERE id = \$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = store.DeleteProject(1)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestProjectStore_GetProjectAssignableUsers_Success проверяет получение
// списка пользователей, которых можно назначить на задачи.
func TestProjectStore_GetProjectAssignableUsers_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewProjectStore(db)
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "email", "full_name", "role", "created_at"}).
		AddRow(1, "a@example.com", "Alice", "user", createdAt).
		AddRow(2, "b@example.com", "Bob", "user", createdAt)

	mock.ExpectQuery(`SELECT DISTINCT u.id, u.email, u.full_name, u.role, u.created_at FROM users u JOIN project_members pm ON pm.user_id = u.id WHERE pm.project_id = \$1 ORDER BY u.full_name`).
		WithArgs(1).
		WillReturnRows(rows)

	users, err := store.GetProjectAssignableUsers(1)

	require.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "Alice", users[0].FullName)
	assert.Equal(t, "Bob", users[1].FullName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

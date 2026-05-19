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

// TestTaskStore_GetTaskByID_Success проверяет получение задачи по ID.
func TestTaskStore_GetTaskByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTaskStore(db)
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "project_id", "title", "description", "status", "priority",
		"assignee_id", "created_by", "due_date", "created_at", "updated_at",
		"type", "hypothesis_id", "resource_id", "conclusion", "task_num", "sprint_id",
		"research_contribution", "research_method",
		"parameters", "metrics", "doi", "tags",
	}).AddRow(1, 1, "Task 1", "Desc", "todo", "high",
		nil, 1, nil, createdAt, nil,
		nil, nil, nil, "", 1, nil,
		"contrib", "method",
		nil, nil, "", "")

	mock.ExpectQuery(`SELECT t.id, t.project_id, t.title, t.description, t.status, t.priority, t.assignee_id, t.created_by, t.due_date, t.created_at, t.updated_at, t.type, t.hypothesis_id, t.resource_id, t.conclusion, t.task_num, t.sprint_id, t.research_contribution, t.research_method, t.parameters, t.metrics, COALESCE\(t.doi, ''\) AS doi, COALESCE\(\(SELECT STRING_AGG\(tg.name, ', '\) FROM tags tg JOIN task_tags tt ON tg.id = tt.tag_id WHERE tt.task_id = t.id\), ''\) as tags FROM tasks t WHERE t.id=\$1`).
		WithArgs(1).
		WillReturnRows(rows)

	task, err := store.GetTaskByID(1)

	require.NoError(t, err)
	assert.Equal(t, 1, task.Id)
	assert.Equal(t, "Task 1", task.Title)
	assert.Equal(t, "todo", task.Status)
	assert.Equal(t, "high", task.Priority)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTaskStore_GetTaskByID_NotFound проверяет обработку отсутствующей задачи.
func TestTaskStore_GetTaskByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTaskStore(db)

	mock.ExpectQuery(`SELECT t.id, t.project_id, t.title, t.description, t.status, t.priority, t.assignee_id, t.created_by, t.due_date, t.created_at, t.updated_at, t.type, t.hypothesis_id, t.resource_id, t.conclusion, t.task_num, t.sprint_id, t.research_contribution, t.research_method, t.parameters, t.metrics, COALESCE\(t.doi, ''\) AS doi, COALESCE\(\(SELECT STRING_AGG\(tg.name, ', '\) FROM tags tg JOIN task_tags tt ON tg.id = tt.tag_id WHERE tt.task_id = t.id\), ''\) as tags FROM tasks t WHERE t.id=\$1`).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	_, err = store.GetTaskByID(999)

	assert.EqualError(t, err, "not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTaskStore_UpdateTaskSprint_Success проверяет обновление спринта задачи.
func TestTaskStore_UpdateTaskSprint_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTaskStore(db)
	sprintID := 5

	mock.ExpectExec(`UPDATE tasks SET sprint_id = \$1, updated_at = NOW\(\) WHERE id = \$2`).
		WithArgs(sprintID, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = store.UpdateTaskSprint(1, &sprintID)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTaskStore_UpdateTaskSprint_NilSprint проверяет сброс спринта (NULL).
func TestTaskStore_UpdateTaskSprint_NilSprint(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTaskStore(db)

	mock.ExpectExec(`UPDATE tasks SET sprint_id = \$1, updated_at = NOW\(\) WHERE id = \$2`).
		WithArgs(nil, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = store.UpdateTaskSprint(1, nil)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTaskStore_GetProjectReportData_Success проверяет генерацию отчета
// по ГОСТ 7.32 для проекта с задачами.
func TestTaskStore_GetProjectReportData_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTaskStore(db)

	mock.ExpectQuery(`SELECT name, description FROM projects WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"name", "description"}).AddRow("AI Project", "AI Description"))

	rows := sqlmock.NewRows([]string{"task_num", "title", "description", "conclusion"}).
		AddRow(1, "Experiment", "Did stuff", "Success")

	mock.ExpectQuery(`SELECT task_num, title, COALESCE\(description, ''\), COALESCE\(conclusion, ''\) FROM tasks WHERE project_id = \$1 AND status IN \('done', 'ГОТОВО'\) ORDER BY task_num`).
		WithArgs(1).
		WillReturnRows(rows)

	report, err := store.GetProjectReportData(1)

	require.NoError(t, err)
	assert.Contains(t, report, "AI Project")
	assert.Contains(t, report, "Experiment")
	assert.Contains(t, report, "Success")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTaskStore_GetProjectReportData_NoTasks проверяет генерацию отчета
// для проекта без завершенных задач.
func TestTaskStore_GetProjectReportData_NoTasks(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTaskStore(db)

	mock.ExpectQuery(`SELECT name, description FROM projects WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"name", "description"}).AddRow("Empty Project", "No desc"))

	mock.ExpectQuery(`SELECT task_num, title, COALESCE\(description, ''\), COALESCE\(conclusion, ''\) FROM tasks WHERE project_id = \$1 AND status IN \('done', 'ГОТОВО'\) ORDER BY task_num`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"task_num", "title", "description", "conclusion"}))

	report, err := store.GetProjectReportData(1)

	require.NoError(t, err)
	assert.Contains(t, report, "Empty Project")
	assert.Contains(t, report, "ЗАКЛЮЧЕНИЕ")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTaskStore_CreateHypothesis_Success проверяет создание гипотезы.
func TestTaskStore_CreateHypothesis_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTaskStore(db)
	createdAt := time.Now()

	mock.ExpectQuery(`INSERT INTO hypotheses \(project_id, title, description, created_by\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id, created_at`).
		WithArgs(1, "H1", "Desc", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(5, createdAt))

	h, err := store.CreateHypothesis(models.Hypothesis{
		ProjectID: 1,
		Title:     "H1",
		Description: "Desc",
		CreatedBy: 1,
	})

	require.NoError(t, err)
	assert.Equal(t, 5, h.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTaskStore_GetProjectHypotheses_Success проверяет получение
// списка гипотез проекта.
func TestTaskStore_GetProjectHypotheses_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTaskStore(db)
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "project_id", "title", "description", "status", "created_by", "created_at"}).
		AddRow(1, 1, "Hyp 1", "Desc 1", "active", 1, createdAt).
		AddRow(2, 1, "Hyp 2", "Desc 2", "rejected", 1, createdAt)

	mock.ExpectQuery(`SELECT id, project_id, title, COALESCE\(description, ''\), status, created_by, created_at FROM hypotheses WHERE project_id = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	list, err := store.GetProjectHypotheses(1)

	require.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, "Hyp 1", list[0].Title)
	assert.Equal(t, "rejected", list[1].Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTaskStore_GetProjectHypotheses_Empty проверяет, что пустой
// результат возвращается как пустой слайс.
func TestTaskStore_GetProjectHypotheses_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTaskStore(db)

	mock.ExpectQuery(`SELECT id, project_id, title, COALESCE\(description, ''\), status, created_by, created_at FROM hypotheses WHERE project_id = \$1`).
		WithArgs(99).
		WillReturnRows(sqlmock.NewRows([]string{"id", "project_id", "title", "description", "status", "created_by", "created_at"}))

	list, err := store.GetProjectHypotheses(99)

	require.NoError(t, err)
	assert.Empty(t, list)
	assert.NotNil(t, list)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTaskStore_UpdateTask_Success проверяет обновление задачи.
func TestTaskStore_UpdateTask_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTaskStore(db)

	mock.ExpectExec(`UPDATE tasks SET title = \$1, description = \$2, priority = \$3, assignee_id = \$4, due_date = \$5, type = \$6, research_contribution = \$7, research_method = \$8, updated_at = NOW\(\) WHERE id = \$9`).
		WithArgs("New Title", "New Desc", "low", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = store.UpdateTask(1, models.Task{
		Title:       "New Title",
		Description: "New Desc",
		Priority:    "LOW",
	})

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTaskStore_DeleteTask_Success проверяет каскадное удаление задачи.
func TestTaskStore_DeleteTask_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTaskStore(db)

	mock.ExpectExec(`DELETE FROM task_tags WHERE task_id = \$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(`DELETE FROM task_history WHERE task_id = \$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(`DELETE FROM comments WHERE entity_type = 'task' AND entity_id = \$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(`DELETE FROM tasks WHERE id = \$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	err = store.DeleteTask(1)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTaskStore_GetTasksByProject_Success проверяет получение
// задач проекта с тегами.
func TestTaskStore_GetTasksByProject_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTaskStore(db)
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "project_id", "title", "description", "status", "priority",
		"assignee_id", "created_by", "due_date", "created_at", "updated_at",
		"type", "hypothesis_id", "resource_id", "conclusion", "task_num", "sprint_id",
		"research_contribution", "research_method",
		"parameters", "metrics", "doi", "tags",
	}).AddRow(1, 1, "T1", "D1", "todo", "medium",
		nil, 1, nil, createdAt, nil,
		nil, nil, nil, "", 1, nil,
		"", "",
		nil, nil, "", "tag1, tag2")

	mock.ExpectQuery(`SELECT t.id, t.project_id, t.title, t.description, t.status, t.priority, t.assignee_id, t.created_by, t.due_date, t.created_at, t.updated_at, t.type, t.hypothesis_id, t.resource_id, t.conclusion, t.task_num, t.sprint_id, t.research_contribution, t.research_method, t.parameters, t.metrics, COALESCE\(t.doi, ''\) AS doi, COALESCE\(\(SELECT STRING_AGG\(tg.name, ', '\) FROM tags tg JOIN task_tags tt ON tg.id = tt.tag_id WHERE tt.task_id = t.id\), ''\) as tags FROM tasks t WHERE t.project_id = \$1 ORDER BY t.task_num ASC`).
		WithArgs(1).
		WillReturnRows(rows)

	tasks, err := store.GetTasksByProject(1)

	require.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, "T1", tasks[0].Title)
	assert.Equal(t, "tag1, tag2", tasks[0].Tags)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTaskStore_GetAllUserTasks_Success проверяет получение задач
// текущего пользователя.
func TestTaskStore_GetAllUserTasks_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTaskStore(db)
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "project_id", "task_num", "sprint_id", "title",
		"description", "status", "priority", "assignee_id",
		"created_by", "due_date", "created_at", "updated_at",
		"type", "hypothesis_id", "resource_id", "conclusion",
		"research_contribution", "research_method",
		"parameters", "metrics", "doi", "tags",
	}).AddRow(1, 1, 1, nil, "Task",
		"Desc", "todo", "high", nil,
		1, nil, createdAt, nil,
		nil, nil, nil, "",
		"", "",
		nil, nil, "", "")

	mock.ExpectQuery(`SELECT t.id, t.project_id, COALESCE\(t.task_num, 0\), t.sprint_id, t.title, COALESCE\(t.description, ''\), t.status, t.priority, t.assignee_id, t.created_by, t.due_date, t.created_at, t.updated_at, t.type, t.hypothesis_id, t.resource_id, COALESCE\(t.conclusion, ''\), t.research_contribution, t.research_method, t.parameters, t.metrics, COALESCE\(t.doi, ''\) AS doi, COALESCE\(\( SELECT STRING_AGG\(tg.name, ','\) FROM tags tg JOIN task_tags tt ON tg.id = tt.tag_id WHERE tt.task_id = t.id \), ''\) as tags FROM tasks t JOIN project_members pm ON t.project_id = pm.project_id WHERE pm.user_id = \$1 ORDER BY t.created_at DESC`).
		WithArgs(1).
		WillReturnRows(rows)

	tasks, err := store.GetAllUserTasks(1)

	require.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, "Task", tasks[0].Title)
	assert.Equal(t, 1, tasks[0].LocalId)
	assert.NoError(t, mock.ExpectationsWereMet())
}

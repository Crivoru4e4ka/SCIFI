package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"project-MVP/db"
	"project-MVP/models"
)

var allowedTaskStatuses = map[string]bool{
	"todo":        true,
	"in_progress": true,
	"review":      true,
	"done":        true,
}

func taskStatusIsValid(status string) bool {
	return allowedTaskStatuses[strings.ToLower(strings.TrimSpace(status))]
}


// GetTasksByProject обёртка над DefaultTaskStore.
func GetTasksByProject(projectId int) ([]models.Task, error) {
	return DefaultTaskStore.GetTasksByProject(projectId)
}

// GetTaskByID обёртка над DefaultTaskStore.
func GetTaskByID(id int) (models.Task, error) {
	return DefaultTaskStore.GetTaskByID(id)
}

// UpdateTaskSprint обёртка над DefaultTaskStore.
func UpdateTaskSprint(taskId int, sprintId *int) error {
	return DefaultTaskStore.UpdateTaskSprint(taskId, sprintId)
}

// GetAllUserTasks обёртка над DefaultTaskStore.
func GetAllUserTasks(userID int) ([]models.Task, error) {
	return DefaultTaskStore.GetAllUserTasks(userID)
}

// GetProjectReportData обёртка над DefaultTaskStore.
func GetProjectReportData(projectID int) (string, error) {
	return DefaultTaskStore.GetProjectReportData(projectID)
}

// CreateHypothesis обёртка над DefaultTaskStore.
func CreateHypothesis(h models.Hypothesis) (models.Hypothesis, error) {
	return DefaultTaskStore.CreateHypothesis(h)
}

// GetProjectHypotheses обёртка над DefaultTaskStore.
func GetProjectHypotheses(projectID int) ([]models.Hypothesis, error) {
	return DefaultTaskStore.GetProjectHypotheses(projectID)
}

// UpdateTask обёртка над DefaultTaskStore.
func UpdateTask(taskId int, task models.Task) error {
	return DefaultTaskStore.UpdateTask(taskId, task)
}

// DeleteTask обёртка над DefaultTaskStore.
func DeleteTask(taskId int) error {
	return DefaultTaskStore.DeleteTask(taskId)
}

func CreateTask(task models.Task) (models.Task, error) {
	task.Title = strings.TrimSpace(task.Title)
	task.Description = strings.TrimSpace(task.Description)
	task.Priority = strings.ToLower(strings.TrimSpace(task.Priority))

	if task.ProjectId <= 0 {
		return models.Task{}, errors.New("project_id is required")
	}
	if task.Title == "" {
		return models.Task{}, errors.New("title is required")
	}
	if task.CreatedBy <= 0 {
		return models.Task{}, errors.New("created_by is required")
	}
	if task.Priority == "" {
		task.Priority = "medium"
	}

	if _, err := GetProjectByID(task.ProjectId); err != nil {
		return models.Task{}, err
	}
	if _, err := GetUserByID(task.CreatedBy); err != nil {
		return models.Task{}, err
	}

	if task.AssigneeId != nil {
		if _, err := GetUserByID(*task.AssigneeId); err != nil {
			return models.Task{}, err
		}
		member, err := IsUserInProject(task.ProjectId, *task.AssigneeId)
		if err != nil {
			return models.Task{}, err
		}
		if !member {
			return models.Task{}, errors.New("assignee is not a member of the project")
		}
	}

	if strings.TrimSpace(task.Status) == "" {
		task.Status = "todo"
	} else {
		task.Status = strings.ToLower(strings.TrimSpace(task.Status))
		if !taskStatusIsValid(task.Status) {
			task.Status = "todo"
		}
	}

	query := `INSERT INTO tasks (project_id, title, description, status, priority, assignee_id, created_by, due_date, type, hypothesis_id, resource_id, conclusion, parameters, metrics, doi, research_contribution, research_method, created_at, updated_at)
              VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,NOW(),NOW()) RETURNING id, task_num`

	var newId int
	var newTaskNum int
	if err := db.DB.QueryRow(query,
		task.ProjectId, task.Title, task.Description, task.Status, task.Priority,
		task.AssigneeId, task.CreatedBy, task.DueDate, task.Type, task.HypothesisId,
		task.ResourceId, task.Conclusion, task.Parameters, task.Metrics, task.DOI,
		task.ResearchContribution, task.ResearchMethod).Scan(&newId, &newTaskNum); err != nil {
		return models.Task{}, err
	}

	LogActivity(task.CreatedBy, task.ProjectId, "task", newId, "created", "Создал научную задачу: "+task.Title)

	task.Id = newId
	task.TaskNum = newTaskNum
	task.CreatedAt = time.Now()

	if strings.TrimSpace(task.Tags) != "" {
		tagNames := strings.Split(task.Tags, ",")
		for _, name := range tagNames {
			name = strings.TrimSpace(name)
			if name != "" {
				_, err := AddTagToTask(task.Id, models.Tag{Name: name})
				if err != nil {
					log.Printf("Ошибка при добавлении тега [%s]: %v", name, err)
				}
			}
		}
	}
	return task, nil
}

func UpdateTaskStatus(taskId int, userId int, newStatus string) error {
	newStatus = strings.ToLower(strings.TrimSpace(newStatus))
	if !taskStatusIsValid(newStatus) {
		return errors.New("invalid task status")
	}

	var oldStatus string
	err := db.DB.QueryRow(`SELECT status FROM tasks WHERE id=$1`, taskId).Scan(&oldStatus)
	if err != nil {
		return err
	}

	_, err = db.DB.Exec(`UPDATE tasks SET status=$1, updated_at=NOW() WHERE id=$2`, newStatus, taskId)
	if err != nil {
		return err
	}

	LogActivity(userId, taskId, "task", taskId, "status_changed", "Сменил статус задачи на "+newStatus)

	queryHistory := `INSERT INTO task_history (task_id, changed_by, field_name, old_value, new_value) VALUES ($1, $2, $3, $4, $5)`
	if _, err := db.DB.Exec(queryHistory, taskId, userId, "status", oldStatus, newStatus); err != nil {
		log.Printf("Error writing to task_history: %v\n", err)
	}

	return nil
}

// scanTaskFromRows сканирует строку результата запроса в модель Task.
// Используется как вспомогательная функция для избежания дублирования кода.
func scanTaskFromRows(rows *sql.Rows) (models.Task, error) {
	var t models.Task
	var assignee, sprint, hypothesis, resource sql.NullInt64
	var dueDate, updatedAt, taskType sql.NullString
	var paramsBytes, metricsBytes []byte

	err := rows.Scan(
		&t.Id, &t.ProjectId, &t.Title, &t.Description, &t.Status, &t.Priority,
		&assignee, &t.CreatedBy, &dueDate, &t.CreatedAt, &updatedAt,
		&taskType, &hypothesis, &resource, &t.Conclusion, &t.TaskNum, &sprint,
		&t.ResearchContribution, &t.ResearchMethod,
		&paramsBytes, &metricsBytes, &t.DOI,
		&t.Tags,
	)
	if err != nil {
		return t, err
	}

	if assignee.Valid {
		val := int(assignee.Int64)
		t.AssigneeId = &val
	}
	if sprint.Valid {
		val := int(sprint.Int64)
		t.SprintId = &val
	}
	if hypothesis.Valid {
		val := int(hypothesis.Int64)
		t.HypothesisId = &val
	}
	if resource.Valid {
		val := int(resource.Int64)
		t.ResourceId = &val
	}
	if dueDate.Valid {
		t.DueDate = &dueDate.String
	}
	if updatedAt.Valid {
		t.UpdatedAt = &updatedAt.String
	}
	if taskType.Valid {
		t.Type = &taskType.String
	}
	if len(paramsBytes) > 0 {
		t.Parameters = json.RawMessage(paramsBytes)
	}
	if len(metricsBytes) > 0 {
		t.Metrics = json.RawMessage(metricsBytes)
	}

	return t, nil
}

// scanTaskFromRow сканирует одну строку в модель Task.
func scanTaskFromRow(row *sql.Row) (models.Task, error) {
	var t models.Task
	var assignee, sprint, hypothesis, resource sql.NullInt64
	var dueDate, updatedAt, taskType sql.NullString
	var paramsBytes, metricsBytes []byte

	err := row.Scan(
		&t.Id, &t.ProjectId, &t.Title, &t.Description, &t.Status, &t.Priority,
		&assignee, &t.CreatedBy, &dueDate, &t.CreatedAt, &updatedAt,
		&taskType, &hypothesis, &resource, &t.Conclusion, &t.TaskNum, &sprint,
		&t.ResearchContribution, &t.ResearchMethod,
		&paramsBytes, &metricsBytes, &t.DOI,
		&t.Tags,
	)
	if err != nil {
		return t, err
	}

	if assignee.Valid {
		val := int(assignee.Int64)
		t.AssigneeId = &val
	}
	if sprint.Valid {
		val := int(sprint.Int64)
		t.SprintId = &val
	}
	if hypothesis.Valid {
		val := int(hypothesis.Int64)
		t.HypothesisId = &val
	}
	if resource.Valid {
		val := int(resource.Int64)
		t.ResourceId = &val
	}
	if dueDate.Valid {
		t.DueDate = &dueDate.String
	}
	if updatedAt.Valid {
		t.UpdatedAt = &updatedAt.String
	}
	if taskType.Valid {
		t.Type = &taskType.String
	}
	if len(paramsBytes) > 0 {
		t.Parameters = json.RawMessage(paramsBytes)
	}
	if len(metricsBytes) > 0 {
		t.Metrics = json.RawMessage(metricsBytes)
	}

	return t, nil
}

// FormatTaskKey формирует ключ задачи по шаблону проекта.
func FormatTaskKey(projectKey string, taskNum int) string {
	return fmt.Sprintf("%s-%d", projectKey, taskNum)
}

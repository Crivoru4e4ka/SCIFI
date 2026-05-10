package services

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"project-MVP/db"
	"project-MVP/models"
)

var allowedTaskStatuses = map[string]bool{
	"todo":        true,
	"in_progress": true,
	"done":        true,
}

func taskStatusIsValid(status string) bool {
	return allowedTaskStatuses[strings.ToLower(strings.TrimSpace(status))]
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

	query := `INSERT INTO tasks (project_id, title, description, status, priority, assignee_id, created_by, due_date, created_at, updated_at)
	          VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW(),NOW()) RETURNING id`
	var newId int
	if err := db.DB.QueryRow(query,
		task.ProjectId, task.Title, task.Description, task.Status, task.Priority,
		task.AssigneeId, task.CreatedBy, task.DueDate).Scan(&newId); err != nil {
		return models.Task{}, err
	}

	task.Id = newId
	task.CreatedAt = time.Now()
	return task, nil
}

// Добавь userId в параметры, чтобы знать, кто сменил статус
func UpdateTaskStatus(taskId int, userId int, newStatus string) error {
	newStatus = strings.ToLower(strings.TrimSpace(newStatus))
	if !taskStatusIsValid(newStatus) {
		return errors.New("invalid task status")
	}

	// 1. Получаем старый статус перед обновлением
	var oldStatus string
	err := db.DB.QueryRow(`SELECT status FROM tasks WHERE id=$1`, taskId).Scan(&oldStatus)
	if err != nil {
		return err
	}

	// 2. Обновляем статус задачи
	_, err = db.DB.Exec(`UPDATE tasks SET status=$1, updated_at=NOW() WHERE id=$2`, newStatus, taskId)
	if err != nil {
		return err
	}

	// 3. АВТОМАТИЧЕСКАЯ ЗАПИСЬ В ИСТОРИЮ (Важно для диплома!)
	queryHistory := `INSERT INTO task_history (task_id, user_id, old_status, new_status) VALUES ($1, $2, $3, $4)`
	_, _ = db.DB.Exec(queryHistory, taskId, userId, oldStatus, newStatus)

	return nil
}

func GetTasksByProject(projectId int) ([]models.Task, error) {
	if _, err := GetProjectByID(projectId); err != nil {
		return nil, err
	}

	rows, err := db.DB.Query(`SELECT id, project_id, title, description, status, priority, assignee_id, created_by, due_date, created_at, updated_at FROM tasks WHERE project_id=$1`, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		var assignee sql.NullInt64
		if err := rows.Scan(&t.Id, &t.ProjectId, &t.Title, &t.Description, &t.Status, &t.Priority, &assignee, &t.CreatedBy, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		if assignee.Valid {
			val := int(assignee.Int64)
			t.AssigneeId = &val
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func GetTaskByID(id int) (models.Task, error) {
	var t models.Task
	var assignee sql.NullInt64
	row := db.DB.QueryRow(`SELECT id, project_id, title, description, status, priority, assignee_id, created_by, due_date, created_at, updated_at FROM tasks WHERE id=$1`, id)
	if err := row.Scan(&t.Id, &t.ProjectId, &t.Title, &t.Description, &t.Status, &t.Priority, &assignee, &t.CreatedBy, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return t, ErrNotFound
		}
		return t, err
	}
	if assignee.Valid {
		val := int(assignee.Int64)
		t.AssigneeId = &val
	}
	return t, nil
}

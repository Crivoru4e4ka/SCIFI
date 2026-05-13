package services

import (
	"database/sql"
	"errors"
	"log"
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
              VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW(),NOW()) RETURNING id, task_num`

	var newId int
	var newTaskNum int
	if err := db.DB.QueryRow(query,
		task.ProjectId, task.Title, task.Description, task.Status, task.Priority,
		task.AssigneeId, task.CreatedBy, task.DueDate).Scan(&newId, &newTaskNum); err != nil {
		return models.Task{}, err
	}

	task.Id = newId
	task.TaskNum = newTaskNum // Теперь здесь будет 1 для первой задачи проекта
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
	queryHistory := `INSERT INTO task_history (task_id, changed_by, field_name, old_value, new_value) VALUES ($1, $2, $3, $4, $5)`
	if _, err := db.DB.Exec(queryHistory, taskId, userId, "status", oldStatus, newStatus); err != nil {
		// Логируем ошибку, но не блокируем обновление статуса
		log.Printf("Error writing to task_history: %v\n", err)
	}

	return nil
}

func GetTasksByProject(projectId int) ([]models.Task, error) {
	query := `
		SELECT 
			id, project_id, title, description, status, priority, 
			assignee_id, created_by, due_date, created_at, updated_at, 
			type, hypothesis_id, resource_id, conclusion, task_num, sprint_id 
		FROM tasks 
		WHERE project_id=$1 
		ORDER BY task_num ASC`

	rows, err := db.DB.Query(query, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		var assignee, sprint, hypothesis, resource sql.NullInt64
		var dueDate, updatedAt, taskType sql.NullString

		// Сканируем все 17 полей
		err := rows.Scan(
			&t.Id, &t.ProjectId, &t.Title, &t.Description, &t.Status, &t.Priority,
			&assignee, &t.CreatedBy, &dueDate, &t.CreatedAt, &updatedAt,
			&taskType, &hypothesis, &resource, &t.Conclusion, &t.TaskNum, &sprint,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			return nil, err
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

		tasks = append(tasks, t)
	}
	return tasks, nil
}

func GetTaskByID(id int) (models.Task, error) {
	var t models.Task
	var assignee, sprint, hypothesis, resource sql.NullInt64
	var dueDate, updatedAt, taskType sql.NullString

	// ПРОВЕРЬ ТУТ: Кавычка должна закрыться сразу после $1
	query := `
		SELECT 
			id, project_id, title, description, status, priority, 
			assignee_id, created_by, due_date, created_at, updated_at, 
			type, hypothesis_id, resource_id, conclusion, task_num, sprint_id 
		FROM tasks 
		WHERE id=$1`

	row := db.DB.QueryRow(query, id)
	err := row.Scan(
		&t.Id, &t.ProjectId, &t.Title, &t.Description, &t.Status, &t.Priority,
		&assignee, &t.CreatedBy, &dueDate, &t.CreatedAt, &updatedAt,
		&taskType, &hypothesis, &resource, &t.Conclusion, &t.TaskNum, &sprint,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return t, ErrNotFound
		}
		return t, err
	}

	// Раскладываем значения
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

	return t, nil
}

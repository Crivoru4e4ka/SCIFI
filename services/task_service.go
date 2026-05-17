package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"project-MVP/db"
	"project-MVP/models"
	"strings"
	"time"
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

	// СОХРАНЕНИЕ ТЕГОВ
	if strings.TrimSpace(task.Tags) != "" {
		tagNames := strings.Split(task.Tags, ",")
		for _, name := range tagNames {
			name = strings.TrimSpace(name)
			if name != "" {
				// ВАЖНО: проверяем ошибку!
				_, err := AddTagToTask(task.Id, models.Tag{Name: name})
				if err != nil {
					log.Printf("Ошибка при добавлении тега [%s]: %v", name, err)
				}
			}
		}
	}
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

	LogActivity(userId, taskId, "task", taskId, "status_changed", "Сменил статус задачи на "+newStatus)

	// 3. АВТОМАТИЧЕСКАЯ ЗАПИСЬ В ИСТОРИЮ (Важно для диплома!)
	queryHistory := `INSERT INTO task_history (task_id, changed_by, field_name, old_value, new_value) VALUES ($1, $2, $3, $4, $5)`
	if _, err := db.DB.Exec(queryHistory, taskId, userId, "status", oldStatus, newStatus); err != nil {
		// Логируем ошибку, но не блокируем обновление статуса
		log.Printf("Error writing to task_history: %v\n", err)
	}

	return nil
}

func GetTasksByProject(projectId int) ([]models.Task, error) {
	// Используем STRING_AGG, чтобы собрать все теги в одну строку "тег1,тег2"
	query := `
		SELECT 
			t.id, t.project_id, t.title, t.description, t.status, t.priority, 
			t.assignee_id, t.created_by, t.due_date, t.created_at, t.updated_at, 
			t.type, t.hypothesis_id, t.resource_id, t.conclusion, t.task_num, t.sprint_id,
			t.research_contribution, t.research_method,
			t.parameters, t.metrics, COALESCE(t.doi, '') AS doi,
			COALESCE((SELECT STRING_AGG(tg.name, ', ') FROM tags tg JOIN task_tags tt ON tg.id = tt.tag_id WHERE tt.task_id = t.id), '') as tags
		FROM tasks t
		WHERE t.project_id = $1
		ORDER BY t.task_num ASC`

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
			log.Printf("Scan error in GetTasksByProject: %v", err)
			continue
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

		tasks = append(tasks, t)
	}
	return tasks, nil
}

func GetTaskByID(id int) (models.Task, error) {
	var t models.Task
	var assignee, sprint, hypothesis, resource sql.NullInt64
	var dueDate, updatedAt, taskType sql.NullString

	// Добавляем получение тегов и в эту функцию тоже!
	query := `
		SELECT 
			t.id, t.project_id, t.title, t.description, t.status, t.priority, 
			t.assignee_id, t.created_by, t.due_date, t.created_at, t.updated_at, 
			t.type, t.hypothesis_id, t.resource_id, t.conclusion, t.task_num, t.sprint_id,
			t.research_contribution, t.research_method,
			t.parameters, t.metrics, COALESCE(t.doi, '') AS doi,
			COALESCE((SELECT STRING_AGG(tg.name, ', ') FROM tags tg JOIN task_tags tt ON tg.id = tt.tag_id WHERE tt.task_id = t.id), '') as tags
		FROM tasks t WHERE t.id=$1`

	var paramsBytes, metricsBytes []byte
	err := db.DB.QueryRow(query, id).Scan(
		&t.Id, &t.ProjectId, &t.Title, &t.Description, &t.Status, &t.Priority,
		&assignee, &t.CreatedBy, &dueDate, &t.CreatedAt, &updatedAt,
		&taskType, &hypothesis, &resource, &t.Conclusion, &t.TaskNum, &sprint,
		&t.ResearchContribution, &t.ResearchMethod,
		&paramsBytes, &metricsBytes, &t.DOI,
		&t.Tags,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return t, errors.New("not found")
		}
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

func UpdateTaskSprint(taskId int, sprintId *int) error {
	// Если sprintId == nil, в базе запишется NULL (задача вернется в бэклог)
	query := `UPDATE tasks SET sprint_id = $1, updated_at = NOW() WHERE id = $2`
	_, err := db.DB.Exec(query, sprintId, taskId)
	return err
}

func GetAllUserTasks(userID int) ([]models.Task, error) {
	query := `
		SELECT 
			t.id, t.project_id, COALESCE(t.task_num, 0), t.sprint_id, t.title, 
			COALESCE(t.description, ''), t.status, t.priority, t.assignee_id, 
			t.created_by, t.due_date, t.created_at, t.updated_at, 
			t.type, t.hypothesis_id, t.resource_id, COALESCE(t.conclusion, ''),
			t.research_contribution, t.research_method,
			t.parameters, t.metrics, COALESCE(t.doi, '') AS doi,
			COALESCE((
				SELECT STRING_AGG(tg.name, ',') 
				FROM tags tg 
				JOIN task_tags tt ON tg.id = tt.tag_id 
				WHERE tt.task_id = t.id
			), '') as tags
		FROM tasks t
		JOIN project_members pm ON t.project_id = pm.project_id
		WHERE pm.user_id = $1
		ORDER BY t.created_at DESC`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		var paramsBytes, metricsBytes []byte
		// Используем правильные имена полей из твоей модели (Id, ProjectId...)
		err := rows.Scan(
			&t.Id, &t.ProjectId, &t.TaskNum, &t.SprintId, &t.Title,
			&t.Description, &t.Status, &t.Priority, &t.AssigneeId,
			&t.CreatedBy, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
			&t.Type, &t.HypothesisId, &t.ResourceId, &t.Conclusion, &t.ResearchContribution, &t.ResearchMethod,
			&paramsBytes, &metricsBytes, &t.DOI,
			&t.Tags,
		)
		if err != nil {
			log.Printf("Scan error in GetAllUserTasks: %v", err)
			continue
		}
		if len(paramsBytes) > 0 {
			t.Parameters = json.RawMessage(paramsBytes)
		}
		if len(metricsBytes) > 0 {
			t.Metrics = json.RawMessage(metricsBytes)
		}
		t.LocalId = t.TaskNum // Синхронизируем для фронтенда
		tasks = append(tasks, t)
	}

	if tasks == nil {
		tasks = []models.Task{}
	}
	return tasks, nil
}

func GetProjectReportData(projectID int) (string, error) {
	var projectName string
	var projectDesc string
	// Получаем данные проекта
	err := db.DB.QueryRow("SELECT name, description FROM projects WHERE id = $1", projectID).Scan(&projectName, &projectDesc)
	if err != nil {
		return "", err
	}

	// Получаем все завершенные задачи
	query := `SELECT task_num, title, COALESCE(description, ''), COALESCE(conclusion, '') 
	          FROM tasks WHERE project_id = $1 AND status IN ('done', 'ГОТОВО') ORDER BY task_num`
	rows, err := db.DB.Query(query, projectID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// Формируем текст по структуре ГОСТ 7.32
	report := "ОТЧЕТ О НАУЧНО-ИССЛЕДОВАТЕЛЬСКОЙ РАБОТЕ\n"
	report += "Тема: " + projectName + "\n\n"
	report += "1. ВВЕДЕНИЕ\n"
	report += projectDesc + "\n\n"
	report += "2. ОСНОВНАЯ ЧАСТЬ (РЕЗУЛЬТАТЫ ЭТАПОВ)\n"

	for rows.Next() {
		var num int
		var title, desc, conc string
		rows.Scan(&num, &title, &desc, &conc)
		report += fmt.Sprintf("\nЭтап %d: %s\n", num, title)
		report += "Описание работ: " + desc + "\n"
		if conc != "" {
			report += "Научный вывод: " + conc + "\n"
		}
	}

	report += "\n\n3. ЗАКЛЮЧЕНИЕ\n"
	report += "Задачи этапа НИР выполнены в полном объеме."

	return report, nil
}

// Создать новую гипотезу
func CreateHypothesis(h models.Hypothesis) (models.Hypothesis, error) {
	query := `INSERT INTO hypotheses (project_id, title, description, created_by) 
              VALUES ($1, $2, $3, $4) RETURNING id, created_at`

	// Используем h.ID и h.ProjectID (заглавными)
	err := db.DB.QueryRow(query, h.ProjectID, h.Title, h.Description, h.CreatedBy).
		Scan(&h.ID, &h.CreatedAt)

	return h, err
}

// Получить все гипотезы проекта
func GetProjectHypotheses(projectID int) ([]models.Hypothesis, error) {
	query := `SELECT id, project_id, title, COALESCE(description, ''), status, created_by, created_at 
              FROM hypotheses WHERE project_id = $1`

	rows, err := db.DB.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Hypothesis
	for rows.Next() {
		var h models.Hypothesis
		// Используем h.ID и h.ProjectID (заглавными)
		err := rows.Scan(&h.ID, &h.ProjectID, &h.Title, &h.Description, &h.Status, &h.CreatedBy, &h.CreatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, h)
	}

	if list == nil {
		list = []models.Hypothesis{}
	}
	return list, nil
}

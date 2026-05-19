package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"project-MVP/db"
	"project-MVP/models"
)

// TaskStore инкапсулирует операции с задачами.
type TaskStore struct {
	DB db.DBPool
}

// NewTaskStore создает новый экземпляр TaskStore.
func NewTaskStore(database db.DBPool) *TaskStore {
	return &TaskStore{DB: database}
}

// GetTasksByProject возвращает все задачи проекта.
func (s *TaskStore) GetTasksByProject(projectId int) ([]models.Task, error) {
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

	rows, err := s.DB.Query(query, projectId)
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

// GetTaskByID возвращает задачу по ID.
func (s *TaskStore) GetTaskByID(id int) (models.Task, error) {
	var t models.Task
	var assignee, sprint, hypothesis, resource sql.NullInt64
	var dueDate, updatedAt, taskType sql.NullString

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
	err := s.DB.QueryRow(query, id).Scan(
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

// UpdateTaskSprint обновляет спринт задачи.
func (s *TaskStore) UpdateTaskSprint(taskId int, sprintId *int) error {
	query := `UPDATE tasks SET sprint_id = $1, updated_at = NOW() WHERE id = $2`
	_, err := s.DB.Exec(query, sprintId, taskId)
	return err
}

// GetAllUserTasks возвращает все задачи пользователя.
func (s *TaskStore) GetAllUserTasks(userID int) ([]models.Task, error) {
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

	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		var paramsBytes, metricsBytes []byte
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
		t.LocalId = t.TaskNum
		tasks = append(tasks, t)
	}

	if tasks == nil {
		tasks = []models.Task{}
	}
	return tasks, nil
}

// GetProjectReportData формирует текст отчета по ГОСТ 7.32.
func (s *TaskStore) GetProjectReportData(projectID int) (string, error) {
	var projectName string
	var projectDesc string
	err := s.DB.QueryRow("SELECT name, description FROM projects WHERE id = $1", projectID).Scan(&projectName, &projectDesc)
	if err != nil {
		return "", err
	}

	query := `SELECT task_num, title, COALESCE(description, ''), COALESCE(conclusion, '') 
	          FROM tasks WHERE project_id = $1 AND status IN ('done', 'ГОТОВО') ORDER BY task_num`
	rows, err := s.DB.Query(query, projectID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

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

// CreateHypothesis создает новую гипотезу.
func (s *TaskStore) CreateHypothesis(h models.Hypothesis) (models.Hypothesis, error) {
	query := `INSERT INTO hypotheses (project_id, title, description, created_by) 
              VALUES ($1, $2, $3, $4) RETURNING id, created_at`

	err := s.DB.QueryRow(query, h.ProjectID, h.Title, h.Description, h.CreatedBy).
		Scan(&h.ID, &h.CreatedAt)

	return h, err
}

// GetProjectHypotheses возвращает все гипотезы проекта.
func (s *TaskStore) GetProjectHypotheses(projectID int) ([]models.Hypothesis, error) {
	query := `SELECT id, project_id, title, COALESCE(description, ''), status, created_by, created_at 
              FROM hypotheses WHERE project_id = $1`

	rows, err := s.DB.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Hypothesis
	for rows.Next() {
		var h models.Hypothesis
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

// UpdateTask обновляет данные задачи.
func (s *TaskStore) UpdateTask(taskId int, task models.Task) error {
	query := `UPDATE tasks SET 
		title = $1, description = $2, priority = $3, assignee_id = $4, 
		due_date = $5, type = $6, research_contribution = $7, research_method = $8,
		updated_at = NOW() 
		WHERE id = $9`

	_, err := s.DB.Exec(query,
		task.Title, task.Description, strings.ToLower(task.Priority),
		task.AssigneeId, task.DueDate, task.Type,
		task.ResearchContribution, task.ResearchMethod, taskId)
	return err
}

// DeleteTask полностью удаляет задачу.
func (s *TaskStore) DeleteTask(taskId int) error {
	_, _ = s.DB.Exec("DELETE FROM task_tags WHERE task_id = $1", taskId)
	_, _ = s.DB.Exec("DELETE FROM task_history WHERE task_id = $1", taskId)
	_, _ = s.DB.Exec("DELETE FROM comments WHERE entity_type = 'task' AND entity_id = $1", taskId)

	_, err := s.DB.Exec("DELETE FROM tasks WHERE id = $1", taskId)
	return err
}

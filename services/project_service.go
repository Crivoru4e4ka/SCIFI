package services

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"project-MVP/db"
	"project-MVP/models"
)

func CreateProject(project models.Project) (models.Project, error) {
	project.Name = strings.TrimSpace(project.Name)
	project.Description = strings.TrimSpace(project.Description)
	project.Status = strings.ToLower(strings.TrimSpace(project.Status))

	if project.Name == "" {
		return models.Project{}, errors.New("project name is required")
	}
	if project.CreatedBy <= 0 {
		return models.Project{}, errors.New("created_by is required")
	}
	if project.StartDate.IsZero() {
		project.StartDate = time.Now()
	}
	if project.Status == "" {
		project.Status = "active"
	}
	// Если фронтенд не прислал статус видимости, ставим по умолчанию
	if project.Visibility == "" {
		project.Visibility = "closed"
	}

	if _, err := GetUserByID(project.CreatedBy); err != nil {
		return models.Project{}, err
	}

	createdAt := time.Now()
	tx, err := db.DB.Begin()
	if err != nil {
		return models.Project{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `INSERT INTO projects (name, key, description, start_date, end_date, status, created_by, created_at, research_goal, main_hypothesis, novelty, expected_result, visibility)
	          VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id`
	var newId int
	if err := tx.QueryRow(query, project.Name, project.Key, project.Description, project.StartDate, project.EndDate, project.Status, project.CreatedBy, createdAt, project.ResearchGoal, project.MainHypothesis, project.Novelty, project.ExpectedResult, project.Visibility).Scan(&newId); err != nil {
		return models.Project{}, err
	}

	creatorRole := "manager"
	memberQuery := `INSERT INTO project_members (project_id, user_id, role) VALUES ($1,$2,$3)`
	if _, err := tx.Exec(memberQuery, newId, project.CreatedBy, creatorRole); err != nil {
		return models.Project{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Project{}, err
	}

	project.Id = newId
	project.CreatedAt = createdAt
	return project, nil
}

// GetProjects возвращает вообще все проекты (обычно для админов)
func GetProjects() ([]models.Project, error) {
	rows, err := db.DB.Query(`SELECT id, name, key, description, research_goal, main_hypothesis, novelty, expected_result, visibility, start_date, end_date, status, created_by, created_at FROM projects`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.Id, &p.Name, &p.Key, &p.Description, &p.ResearchGoal, &p.MainHypothesis, &p.Novelty, &p.ExpectedResult, &p.Visibility, &p.StartDate, &p.EndDate, &p.Status, &p.CreatedBy, &p.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

// GetUserProjects — САМАЯ ВАЖНАЯ ФУНКЦИЯ. Реализует доступ: Мои + Где я участник + Открытые
func GetUserProjects(userId int) ([]models.Project, error) {
	query := `
		SELECT DISTINCT p.id, p.name, p.key, p.description, p.research_goal, 
		       p.main_hypothesis, p.novelty, p.expected_result, p.visibility,
		       p.start_date, p.end_date, p.status, p.created_by, p.created_at 
		FROM projects p
		LEFT JOIN project_members pm ON p.id = pm.project_id
		WHERE p.created_by = $1 
		   OR pm.user_id = $1 
		   OR p.visibility = 'open'
		ORDER BY p.created_at DESC`

	rows, err := db.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.Id, &p.Name, &p.Key, &p.Description, &p.ResearchGoal, &p.MainHypothesis, &p.Novelty, &p.ExpectedResult, &p.Visibility, &p.StartDate, &p.EndDate, &p.Status, &p.CreatedBy, &p.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}

	if list == nil {
		list = []models.Project{}
	}
	return list, nil
}

func GetProjectByID(id int) (models.Project, error) {
	var p models.Project
	query := `SELECT id, name, key, description, research_goal, main_hypothesis, novelty, expected_result, visibility, start_date, end_date, status, created_by, created_at 
	          FROM projects WHERE id=$1`
	row := db.DB.QueryRow(query, id)
	if err := row.Scan(&p.Id, &p.Name, &p.Key, &p.Description, &p.ResearchGoal, &p.MainHypothesis, &p.Novelty, &p.ExpectedResult, &p.Visibility, &p.StartDate, &p.EndDate, &p.Status, &p.CreatedBy, &p.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return p, ErrNotFound
		}
		return p, err
	}
	return p, nil
}

func GetProjectByName(name string) (models.Project, error) {
	var p models.Project
	query := `SELECT id, name, key, description, research_goal, main_hypothesis, novelty, expected_result, visibility, start_date, end_date, status, created_by, created_at 
	          FROM projects WHERE name=$1`
	row := db.DB.QueryRow(query, name)
	if err := row.Scan(&p.Id, &p.Name, &p.Key, &p.Description, &p.ResearchGoal, &p.MainHypothesis, &p.Novelty, &p.ExpectedResult, &p.Visibility, &p.StartDate, &p.EndDate, &p.Status, &p.CreatedBy, &p.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return p, ErrNotFound
		}
		return p, err
	}
	return p, nil
}

func GetProjectProgress(projectId int) (float64, error) {
	if _, err := GetProjectByID(projectId); err != nil {
		return 0, err
	}

	var doneCount, totalCount int
	query := `SELECT
		COUNT(CASE WHEN status = 'done' THEN 1 END) AS done_count,
		COUNT(*) AS total_count
		FROM tasks WHERE project_id=$1`
	row := db.DB.QueryRow(query, projectId)
	if err := row.Scan(&doneCount, &totalCount); err != nil {
		return 0, err
	}
	if totalCount == 0 {
		return 0, nil
	}
	return float64(doneCount) * 100.0 / float64(totalCount), nil
}

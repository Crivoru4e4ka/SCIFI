package services

import (
	"database/sql"
	"errors"
	"time"

	"project-MVP/db"
	"project-MVP/models"
)

// ProjectStore инкапсулирует операции с проектами.
// Принимает db.DBPool, что позволяет мокать зависимость в unit-тестах.
type ProjectStore struct {
	DB db.DBPool
}

// NewProjectStore создает новый экземпляр ProjectStore.
func NewProjectStore(database db.DBPool) *ProjectStore {
	return &ProjectStore{DB: database}
}

// GetProjects возвращает вообще все проекты (обычно для админов).
func (s *ProjectStore) GetProjects() ([]models.Project, error) {
	rows, err := s.DB.Query(`SELECT id, name, key, description, research_goal, main_hypothesis, novelty, expected_result, visibility, start_date, end_date, status, created_by, created_at, team_id, execution_type FROM projects`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Project
	for rows.Next() {
		var p models.Project
		var teamID sql.NullInt64
		if err := rows.Scan(&p.Id, &p.Name, &p.Key, &p.Description, &p.ResearchGoal, &p.MainHypothesis, &p.Novelty, &p.ExpectedResult, &p.Visibility, &p.StartDate, &p.EndDate, &p.Status, &p.CreatedBy, &p.CreatedAt, &teamID, &p.ExecutionType); err != nil {
			return nil, err
		}
		if teamID.Valid {
			v := int(teamID.Int64)
			p.TeamId = &v
		}
		list = append(list, p)
	}
	return list, nil
}

// GetUserProjects — САМАЯ ВАЖНАЯ ФУНКЦИЯ. Реализует доступ: Мои + Где я участник + Открытые + Команда.
func (s *ProjectStore) GetUserProjects(userId int) ([]models.Project, error) {
	query := `
		SELECT DISTINCT p.id, p.name, p.key, p.description, p.research_goal, 
		       p.main_hypothesis, p.novelty, p.expected_result, p.visibility,
		       p.start_date, p.end_date, p.status, p.created_by, p.created_at,
		       p.team_id, p.execution_type
		FROM projects p
		LEFT JOIN project_members pm ON p.id = pm.project_id
		WHERE p.created_by = $1 
		   OR pm.user_id = $1 
		   OR p.visibility = 'open'
		ORDER BY p.created_at DESC`

	rows, err := s.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Project
	for rows.Next() {
		var p models.Project
		var teamID sql.NullInt64
		if err := rows.Scan(&p.Id, &p.Name, &p.Key, &p.Description, &p.ResearchGoal, &p.MainHypothesis, &p.Novelty, &p.ExpectedResult, &p.Visibility, &p.StartDate, &p.EndDate, &p.Status, &p.CreatedBy, &p.CreatedAt, &teamID, &p.ExecutionType); err != nil {
			return nil, err
		}
		if teamID.Valid {
			v := int(teamID.Int64)
			p.TeamId = &v
		}
		list = append(list, p)
	}

	if list == nil {
		list = []models.Project{}
	}
	return list, nil
}

// GetProjectByID возвращает проект по ID.
func (s *ProjectStore) GetProjectByID(id int) (models.Project, error) {
	var p models.Project
	query := `SELECT id, name, key, description, research_goal, main_hypothesis, novelty, expected_result, visibility, start_date, end_date, status, created_by, created_at, team_id, execution_type
	          FROM projects WHERE id=$1`
	row := s.DB.QueryRow(query, id)
	var teamID sql.NullInt64
	if err := row.Scan(&p.Id, &p.Name, &p.Key, &p.Description, &p.ResearchGoal, &p.MainHypothesis, &p.Novelty, &p.ExpectedResult, &p.Visibility, &p.StartDate, &p.EndDate, &p.Status, &p.CreatedBy, &p.CreatedAt, &teamID, &p.ExecutionType); err != nil {
		if err == sql.ErrNoRows {
			return p, ErrNotFound
		}
		return p, err
	}
	if teamID.Valid {
		v := int(teamID.Int64)
		p.TeamId = &v
	}
	return p, nil
}

// GetProjectByName возвращает проект по названию.
func (s *ProjectStore) GetProjectByName(name string) (models.Project, error) {
	var p models.Project
	query := `SELECT id, name, key, description, research_goal, main_hypothesis, novelty, expected_result, visibility, start_date, end_date, status, created_by, created_at, team_id, execution_type
	          FROM projects WHERE name=$1`
	row := s.DB.QueryRow(query, name)
	var teamID sql.NullInt64
	if err := row.Scan(&p.Id, &p.Name, &p.Key, &p.Description, &p.ResearchGoal, &p.MainHypothesis, &p.Novelty, &p.ExpectedResult, &p.Visibility, &p.StartDate, &p.EndDate, &p.Status, &p.CreatedBy, &p.CreatedAt, &teamID, &p.ExecutionType); err != nil {
		if err == sql.ErrNoRows {
			return p, ErrNotFound
		}
		return p, err
	}
	if teamID.Valid {
		v := int(teamID.Int64)
		p.TeamId = &v
	}
	return p, nil
}

// GetProjectProgress возвращает процент завершенных задач в проекте.
func (s *ProjectStore) GetProjectProgress(projectId int) (float64, error) {
	if _, err := s.GetProjectByID(projectId); err != nil {
		return 0, err
	}

	var doneCount, totalCount int
	query := `SELECT
		COUNT(CASE WHEN status = 'done' THEN 1 END) AS done_count,
		COUNT(*) AS total_count
		FROM tasks WHERE project_id=$1`
	row := s.DB.QueryRow(query, projectId)
	if err := row.Scan(&doneCount, &totalCount); err != nil {
		return 0, err
	}
	if totalCount == 0 {
		return 0, nil
	}
	return float64(doneCount) * 100.0 / float64(totalCount), nil
}

// GetProjectAssignableUsers возвращает пользователей, которых можно назначить исполнителями задач в проекте.
func (s *ProjectStore) GetProjectAssignableUsers(projectID int) ([]models.User, error) {
	query := `
		SELECT DISTINCT u.id, u.email, u.full_name, u.role, u.created_at
		FROM users u
		JOIN project_members pm ON pm.user_id = u.id
		WHERE pm.project_id = $1
		ORDER BY u.full_name`

	rows, err := s.DB.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.Id, &u.Email, &u.FullName, &u.Role, &u.CreatedAt); err != nil {
			continue
		}
		users = append(users, u)
	}
	return users, nil
}

// UpdateProject обновляет данные проекта.
func (s *ProjectStore) UpdateProject(id int, p models.Project) error {
	query := `UPDATE projects SET 
		name = $1, description = $2, research_goal = $3, 
		main_hypothesis = $4, novelty = $5, expected_result = $6,
		status = $7, end_date = $8, visibility = $9
		WHERE id = $10`
	_, err := s.DB.Exec(query, p.Name, p.Description, p.ResearchGoal,
		p.MainHypothesis, p.Novelty, p.ExpectedResult,
		p.Status, p.EndDate, p.Visibility, id)
	return err
}

// DeleteProject удаляет проект со всеми связанными данными.
func (s *ProjectStore) DeleteProject(id int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tx.Exec("DELETE FROM tasks WHERE project_id = $1", id)
	tx.Exec("DELETE FROM project_members WHERE project_id = $1", id)
	tx.Exec("DELETE FROM audit_log WHERE project_id = $1", id)
	tx.Exec("DELETE FROM projects WHERE id = $1", id)

	return tx.Commit()
}

// CreateProject создает новый проект с валидацией и транзакцией.
func (s *ProjectStore) CreateProject(project models.Project) (models.Project, error) {
	if err := ValidateAndNormalizeProject(&project); err != nil {
		return models.Project{}, err
	}

	if _, err := GetUserByID(project.CreatedBy); err != nil {
		return models.Project{}, err
	}

	leadRoleID, err := getRoleIDByName("project_lead")
	if err != nil {
		InitRoles()
		leadRoleID, err = getRoleIDByName("project_lead")
		if err != nil {
			return models.Project{}, errors.New("project_lead role not found in RBAC system")
		}
	}

	createdAt := time.Now()
	tx, err := s.DB.Begin()
	if err != nil {
		return models.Project{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `INSERT INTO projects (name, key, description, start_date, end_date, status, created_by, created_at, research_goal, main_hypothesis, novelty, expected_result, visibility, team_id, execution_type)
	          VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) RETURNING id`
	var newId int
	if err := tx.QueryRow(query, project.Name, project.Key, project.Description, project.StartDate, project.EndDate, project.Status, project.CreatedBy, createdAt,
		project.ResearchGoal, project.MainHypothesis, project.Novelty, project.ExpectedResult, project.Visibility, project.TeamId, project.ExecutionType).Scan(&newId); err != nil {
		return models.Project{}, err
	}

	creatorRole := "project_lead"
	memberQuery := `INSERT INTO project_members (project_id, user_id, role, role_id) VALUES ($1,$2,$3,$4)`
	if _, err := tx.Exec(memberQuery, newId, project.CreatedBy, creatorRole, leadRoleID); err != nil {
		return models.Project{}, err
	}

	if project.ExecutionType == "team" && project.TeamId != nil {
		teamMembers, err := getTeamMembersTx(tx, *project.TeamId)
		if err != nil {
			return models.Project{}, err
		}
		for _, tm := range teamMembers {
			if tm.UserID == project.CreatedBy {
				continue
			}
			memberRole, memberRoleID := resolveTeamMemberRole(tm.Role)
			_, err := tx.Exec(`INSERT INTO project_members (project_id, user_id, role, role_id)
				VALUES ($1,$2,$3,$4)
				ON CONFLICT (project_id, user_id) DO NOTHING`,
				newId, tm.UserID, memberRole, memberRoleID)
			if err != nil {
				return models.Project{}, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return models.Project{}, err
	}

	project.Id = newId
	project.CreatedAt = createdAt
	return project, nil
}

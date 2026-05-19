package services

import (
	"database/sql"

	"project-MVP/db"
	"project-MVP/models"
)

// SprintStore инкапсулирует операции со спринтами.
type SprintStore struct {
	DB db.DBPool
}

// NewSprintStore создает новый экземпляр SprintStore.
func NewSprintStore(database db.DBPool) *SprintStore {
	return &SprintStore{DB: database}
}

// GetProjectSprints возвращает спринты проекта.
func (s *SprintStore) GetProjectSprints(projectId int) ([]models.Sprint, error) {
	rows, err := s.DB.Query(`
		SELECT id, project_id, name, status, start_date, end_date 
		FROM sprints 
		WHERE project_id = $1 
		ORDER BY created_at DESC`,
		projectId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sprints []models.Sprint
	for rows.Next() {
		var sp models.Sprint
		err := rows.Scan(
			&sp.ID,
			&sp.ProjectID,
			&sp.Name,
			&sp.Status,
			&sp.StartDate,
			&sp.EndDate,
		)
		if err != nil {
			return nil, err
		}
		sprints = append(sprints, sp)
	}

	if sprints == nil {
		sprints = []models.Sprint{}
	}
	return sprints, nil
}

// CreateSprint создает новый спринт.
func (s *SprintStore) CreateSprint(sp *models.Sprint) error {
	query := `
        INSERT INTO sprints (project_id, name, status, start_date, end_date) 
        VALUES ($1, $2, $3, $4, $5) 
        RETURNING id`

	err := s.DB.QueryRow(query,
		sp.ProjectID,
		sp.Name,
		sp.Status,
		sp.StartDate,
		sp.EndDate,
	).Scan(&sp.ID)

	return err
}

// GetSprintByID возвращает спринт по ID.
func (s *SprintStore) GetSprintByID(id int) (models.Sprint, error) {
	var sp models.Sprint
	query := `SELECT id, project_id, name, status, start_date, end_date, goal FROM sprints WHERE id = $1`
	err := s.DB.QueryRow(query, id).Scan(&sp.ID, &sp.ProjectID, &sp.Name, &sp.Status, &sp.StartDate, &sp.EndDate, &sp.Goal)
	if err != nil {
		if err == sql.ErrNoRows {
			return sp, ErrNotFound
		}
		return sp, err
	}
	return sp, nil
}

// StartSprint активирует спринт.
func (s *SprintStore) StartSprint(id int, sp models.Sprint) error {
	query := `UPDATE sprints SET name=$1, start_date=$2, end_date=$3, goal=$4, status='active' WHERE id=$5`
	_, err := s.DB.Exec(query, sp.Name, sp.StartDate, sp.EndDate, sp.Goal, id)
	return err
}

// CompleteSprint завершает спринт и возвращает незавершенные задачи в бэклог.
func (s *SprintStore) CompleteSprint(sprintID int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE sprints SET status = 'closed', end_date = NOW() WHERE id = $1`, sprintID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE tasks SET sprint_id = NULL WHERE sprint_id = $1 AND status != 'done'`, sprintID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

package services

import (
	"project-MVP/db"
	"project-MVP/models"
)

// GetProjectSprints возвращает список всех спринтов для конкретного проекта
func GetProjectSprints(projectId int) ([]models.Sprint, error) {
	// Выполняем запрос к БД
	rows, err := db.DB.Query(`
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
		var s models.Sprint
		// Scan автоматически разложит данные по полям структуры.
		// Так как StartDate и EndDate в модели — это *string (указатели),
		// Go сам запишет туда nil, если в базе стоит NULL.
		err := rows.Scan(
			&s.ID,
			&s.ProjectID,
			&s.Name,
			&s.Status,
			&s.StartDate,
			&s.EndDate,
		)
		if err != nil {
			return nil, err
		}
		sprints = append(sprints, s)
	}

	// Если спринтов нет, возвращаем пустой массив вместо nil для красоты JSON
	if sprints == nil {
		sprints = []models.Sprint{}
	}

	return sprints, nil
}

// CreateSprint создает новый спринт в базе данных
func CreateSprint(s *models.Sprint) error {
	query := `
        INSERT INTO sprints (project_id, name, status, start_date, end_date) 
        VALUES ($1, $2, $3, $4, $5) 
        RETURNING id`

	// Выполняем запрос и сразу получаем ID созданной записи
	err := db.DB.QueryRow(query,
		s.ProjectID,
		s.Name,
		s.Status,
		s.StartDate,
		s.EndDate,
	).Scan(&s.ID)

	return err
}

func StartSprint(id int, s models.Sprint) error {
	query := `UPDATE sprints SET name=$1, start_date=$2, end_date=$3, goal=$4, status='active' WHERE id=$5`
	_, err := db.DB.Exec(query, s.Name, s.StartDate, s.EndDate, s.Goal, id)
	return err
}

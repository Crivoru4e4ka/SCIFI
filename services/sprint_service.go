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

func CompleteSprint(sprintID int) error {
	// Начинаем транзакцию
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Закрываем сам спринт
	_, err = tx.Exec(`UPDATE sprints SET status = 'closed', end_date = NOW() WHERE id = $1`, sprintID)
	if err != nil {
		return err
	}

	// 2. Jira-логика: Все задачи, которые НЕ в статусе 'done', возвращаем в бэклог (sprint_id = NULL)
	_, err = tx.Exec(`UPDATE tasks SET sprint_id = NULL WHERE sprint_id = $1 AND status != 'done'`, sprintID)
	if err != nil {
		return err
	}

	// Фиксируем изменения в базе
	return tx.Commit()
}

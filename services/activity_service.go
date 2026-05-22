package services

import (
	"log"
	"project-MVP/db"
	"project-MVP/models"
)

// LogActivity записывает новое действие в таблицу activities.
// Вызывается из других сервисов (TaskService, ProjectService и т.д.)
func LogActivity(userID int, projectID int, entityType string, entityID int, action string, details string) {
	query := `INSERT INTO activities (user_id, project_id, entity_type, entity_id, action, details, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6, NOW())`

	_, err := db.DB.Exec(query, userID, projectID, entityType, entityID, action, details)
	if err != nil {
		log.Printf("Ошибка при записи активности: %v", err)
	}
}

// GetUserActivities возвращает ленту событий для пользователя.
// Показывает действия во всех проектах, в которых состоит пользователь.
func GetUserActivities(userID int) ([]models.Activity, error) {
	query := `
		SELECT 
			a.id, a.user_id, u.full_name, a.project_id, p.name, 
			a.entity_type, a.entity_id, a.action, a.details, a.created_at
		FROM activities a
		JOIN users u ON a.user_id = u.id
		JOIN projects p ON a.project_id = p.id
		JOIN project_members pm ON p.id = pm.project_id
		WHERE pm.user_id = $1
		ORDER BY a.created_at DESC
		LIMIT 50`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activities := []models.Activity{} // Инициализируем пустым слайсом, чтобы не было null в JSON
	for rows.Next() {
		var a models.Activity
		err := rows.Scan(
			&a.Id, &a.UserId, &a.UserName, &a.ProjectId, &a.ProjectName,
			&a.EntityType, &a.EntityId, &a.Action, &a.Details, &a.CreatedAt,
		)
		if err != nil {
			log.Printf("Scan error in GetUserActivities: %v", err)
			continue
		}
		activities = append(activities, a)
	}

	return activities, nil
}

// GetProjectActivities возвращает ленту событий для конкретного проекта.
func GetProjectActivities(projectID int) ([]models.Activity, error) {
	query := `
		SELECT 
			a.id, a.user_id, u.full_name, a.project_id, p.name, 
			a.entity_type, a.entity_id, a.action, a.details, a.created_at
		FROM activities a
		JOIN users u ON a.user_id = u.id
		JOIN projects p ON a.project_id = p.id
		WHERE a.project_id = $1
		ORDER BY a.created_at DESC
		LIMIT 50`

	rows, err := db.DB.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activities := []models.Activity{}
	for rows.Next() {
		var a models.Activity
		err := rows.Scan(
			&a.Id, &a.UserId, &a.UserName, &a.ProjectId, &a.ProjectName,
			&a.EntityType, &a.EntityId, &a.Action, &a.Details, &a.CreatedAt,
		)
		if err != nil {
			log.Printf("Scan error in GetProjectActivities: %v", err)
			continue
		}
		activities = append(activities, a)
	}

	return activities, nil
}

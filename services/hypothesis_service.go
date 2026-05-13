package services

import (
	"database/sql"
	"project-MVP/db"
	"project-MVP/models"
)

// Получить все гипотезы
func GetAllHypotheses() ([]models.Hypothesis, error) {
	rows, err := db.DB.Query(`SELECT id, project_id, title, description, status, created_by, created_at FROM hypotheses`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Hypothesis
	for rows.Next() {
		var h models.Hypothesis
		if err := rows.Scan(&h.ID, &h.ProjectID, &h.Title, &h.Description, &h.Status, &h.CreatedBy, &h.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, h)
	}
	return result, nil
}

// Создать гипотезу
func CreateHypothesis(h models.Hypothesis) (models.Hypothesis, error) {
	row := db.DB.QueryRow(`INSERT INTO hypotheses (project_id, title, description, status, created_by) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`,
		h.ProjectID, h.Title, h.Description, h.Status, h.CreatedBy)
	var id int
	var createdAt sql.NullTime
	if err := row.Scan(&id, &createdAt); err != nil {
		return h, err
	}
	h.ID = id
	if createdAt.Valid {
		h.CreatedAt = createdAt.Time
	}
	return h, nil
}

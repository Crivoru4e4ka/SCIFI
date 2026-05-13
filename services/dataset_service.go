package services

import (
	"database/sql"
	"project-MVP/db"
	"project-MVP/models"
)

// Получить все датасеты
func GetAllDatasets() ([]models.Dataset, error) {
	rows, err := db.DB.Query(`SELECT id, project_id, name, description, version, data_url, parameters, created_at FROM datasets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Dataset
	for rows.Next() {
		var ds models.Dataset
		var parameters []byte
		if err := rows.Scan(&ds.ID, &ds.ProjectID, &ds.Name, &ds.Description, &ds.Version, &ds.DataURL, &parameters, &ds.CreatedAt); err != nil {
			return nil, err
		}
		ds.Parameters = parameters
		result = append(result, ds)
	}
	return result, nil
}

// Создать датасет
func CreateDataset(ds models.Dataset) (models.Dataset, error) {
	row := db.DB.QueryRow(`INSERT INTO datasets (project_id, name, description, version, data_url, parameters) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`,
		ds.ProjectID, ds.Name, ds.Description, ds.Version, ds.DataURL, ds.Parameters)
	var id int
	var createdAt sql.NullTime
	if err := row.Scan(&id, &createdAt); err != nil {
		return ds, err
	}
	ds.ID = id
	if createdAt.Valid {
		ds.CreatedAt = createdAt.Time
	}
	return ds, nil
}

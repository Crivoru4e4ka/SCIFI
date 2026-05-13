package services

import (
	"project-MVP/db"
	"project-MVP/models"
)

// Получить все связи
func GetAllExperimentDatasets() ([]models.ExperimentDataset, error) {
	rows, err := db.DB.Query(`SELECT id, experiment_id, dataset_id FROM experiment_datasets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.ExperimentDataset
	for rows.Next() {
		var eds models.ExperimentDataset
		// Поля ID, ExperimentID, DatasetID должны быть с Большой буквы
		if err := rows.Scan(&eds.ID, &eds.ExperimentID, &eds.DatasetID); err != nil {
			return nil, err
		}
		result = append(result, eds)
	}
	return result, nil
}

// Создать связь
func CreateExperimentDataset(eds models.ExperimentDataset) (models.ExperimentDataset, error) {
	row := db.DB.QueryRow(`INSERT INTO experiment_datasets (experiment_id, dataset_id) VALUES ($1, $2) RETURNING id`,
		eds.ExperimentID, eds.DatasetID)

	var id int
	if err := row.Scan(&id); err != nil {
		return eds, err
	}

	eds.ID = id // Теперь это поле будет доступно
	return eds, nil
}

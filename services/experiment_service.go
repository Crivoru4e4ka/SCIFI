package services

import (
	"project-MVP/db"
	"project-MVP/models"
)

// GetAllExperimentDatasets получает все связи задач с датасетами
func GetAllExperimentDatasets() ([]models.ExperimentDataset, error) {
	rows, err := db.DB.Query(`SELECT task_id, dataset_id FROM experiment_datasets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.ExperimentDataset
	for rows.Next() {
		var eds models.ExperimentDataset
		if err := rows.Scan(&eds.TaskID, &eds.DatasetID); err != nil {
			return nil, err
		}
		result = append(result, eds)
	}
	return result, nil
}

// CreateExperimentDataset создает связь задачи и датасета
func CreateExperimentDataset(eds models.ExperimentDataset) (models.ExperimentDataset, error) {
	_, err := db.DB.Exec(`INSERT INTO experiment_datasets (task_id, dataset_id) VALUES ($1, $2)`,
		eds.TaskID, eds.DatasetID)
	if err != nil {
		return eds, err
	}
	return eds, nil
}

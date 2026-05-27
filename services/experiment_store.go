package services

import (
	"project-MVP/db"
	"project-MVP/models"
)

// ExperimentStore инкапсулирует операции со связями задач и датасетов.
type ExperimentStore struct {
	DB db.DBPool
}

// NewExperimentStore создает новый экземпляр ExperimentStore.
func NewExperimentStore(database db.DBPool) *ExperimentStore {
	return &ExperimentStore{DB: database}
}

// GetAllExperimentDatasets возвращает все связи задач с датасетами.
func (s *ExperimentStore) GetAllExperimentDatasets() ([]models.ExperimentDataset, error) {
	rows, err := s.DB.Query(`SELECT task_id, dataset_id, relation_type FROM experiment_datasets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.ExperimentDataset
	for rows.Next() {
		var eds models.ExperimentDataset
		if err := rows.Scan(&eds.TaskID, &eds.DatasetID, &eds.RelationType); err != nil {
			return nil, err
		}
		result = append(result, eds)
	}
	return result, nil
}

// GetDatasetsByTask возвращает связи задачи с датасетами.
func (s *ExperimentStore) GetDatasetsByTask(taskID int) ([]models.ExperimentDataset, error) {
	rows, err := s.DB.Query(`SELECT task_id, dataset_id, relation_type FROM experiment_datasets WHERE task_id = $1`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.ExperimentDataset
	for rows.Next() {
		var eds models.ExperimentDataset
		if err := rows.Scan(&eds.TaskID, &eds.DatasetID, &eds.RelationType); err != nil {
			return nil, err
		}
		result = append(result, eds)
	}
	return result, nil
}

// CreateExperimentDataset создает связь задачи и датасета.
func (s *ExperimentStore) CreateExperimentDataset(eds models.ExperimentDataset) (models.ExperimentDataset, error) {
	_, err := s.DB.Exec(`INSERT INTO experiment_datasets (task_id, dataset_id, relation_type) VALUES ($1, $2, $3)`,
		eds.TaskID, eds.DatasetID, eds.RelationType)
	if err != nil {
		return eds, err
	}
	return eds, nil
}

// DeleteExperimentDataset удаляет связь задачи и датасета.
func (s *ExperimentStore) DeleteExperimentDataset(taskID, datasetID int) error {
	_, err := s.DB.Exec(`DELETE FROM experiment_datasets WHERE task_id = $1 AND dataset_id = $2`, taskID, datasetID)
	return err
}

// DeleteTaskDatasets удаляет все связи задачи с датасетами.
func (s *ExperimentStore) DeleteTaskDatasets(taskID int) error {
	_, err := s.DB.Exec(`DELETE FROM experiment_datasets WHERE task_id = $1`, taskID)
	return err
}

// GetTasksByDataset возвращает задачи, связанные с датасетом (с relation_type).
func (s *ExperimentStore) GetTasksByDataset(datasetID int) ([]map[string]interface{}, error) {
	rows, err := s.DB.Query(`
		SELECT t.id, t.title, t.type, ed.relation_type
		FROM experiment_datasets ed
		JOIN tasks t ON t.id = ed.task_id
		WHERE ed.dataset_id = $1
		ORDER BY t.id DESC`, datasetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id int
		var title, taskType, relation string
		if err := rows.Scan(&id, &title, &taskType, &relation); err != nil {
			continue
		}
		result = append(result, map[string]interface{}{
			"id":            id,
			"title":         title,
			"type":          taskType,
			"relation_type": relation,
		})
	}
	return result, nil
}

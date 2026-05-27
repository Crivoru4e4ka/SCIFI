package services

import (
	"project-MVP/db"
	"project-MVP/models"
)

// DatasetDependencyStore инкапсулирует операции с lineage.
type DatasetDependencyStore struct {
	DB db.DBPool
}

// NewDatasetDependencyStore создает новый экземпляр.
func NewDatasetDependencyStore(database db.DBPool) *DatasetDependencyStore {
	return &DatasetDependencyStore{DB: database}
}

// CreateDependency создает запись lineage.
func (s *DatasetDependencyStore) CreateDependency(dep models.DatasetDependency) error {
	_, err := s.DB.Exec(
		`INSERT INTO dataset_dependencies (source_dataset_id, target_dataset_id, task_id) VALUES ($1, $2, $3)
		 ON CONFLICT (source_dataset_id, target_dataset_id, task_id) DO NOTHING`,
		dep.SourceDatasetID, dep.TargetDatasetID, dep.TaskID)
	return err
}

// GetLineage возвращает прямые зависимости (source -> target) для датасета.
func (s *DatasetDependencyStore) GetLineage(datasetID int) (map[string]interface{}, error) {
	// Upstream: кто является источником для данного датасета
	upstreamRows, err := s.DB.Query(
		`SELECT d.id, d.name, d.version, dd.task_id FROM dataset_dependencies dd
		 JOIN datasets d ON dd.source_dataset_id = d.id
		 WHERE dd.target_dataset_id = $1`, datasetID)
	if err != nil {
		return nil, err
	}
	defer upstreamRows.Close()

	var upstream []map[string]interface{}
	for upstreamRows.Next() {
		var id, taskID int
		var name, version string
		if err := upstreamRows.Scan(&id, &name, &version, &taskID); err != nil {
			continue
		}
		upstream = append(upstream, map[string]interface{}{
			"id": id, "name": name, "version": version, "task_id": taskID,
		})
	}

	// Downstream: кто зависит от данного датасета
	downstreamRows, err := s.DB.Query(
		`SELECT d.id, d.name, d.version, dd.task_id FROM dataset_dependencies dd
		 JOIN datasets d ON dd.target_dataset_id = d.id
		 WHERE dd.source_dataset_id = $1`, datasetID)
	if err != nil {
		return nil, err
	}
	defer downstreamRows.Close()

	var downstream []map[string]interface{}
	for downstreamRows.Next() {
		var id, taskID int
		var name, version string
		if err := downstreamRows.Scan(&id, &name, &version, &taskID); err != nil {
			continue
		}
		downstream = append(downstream, map[string]interface{}{
			"id": id, "name": name, "version": version, "task_id": taskID,
		})
	}

	return map[string]interface{}{
		"upstream":   upstream,
		"downstream": downstream,
	}, nil
}

// GetFullLineageGraph возвращает полный граф lineage через рекурсивный CTE.
func (s *DatasetDependencyStore) GetFullLineageGraph(datasetID int) (map[string]interface{}, error) {
	// Находим все предки (upstream)
	upRows, err := s.DB.Query(`
		WITH RECURSIVE upstream AS (
			SELECT source_dataset_id, target_dataset_id, task_id, 1 AS depth
			FROM dataset_dependencies
			WHERE target_dataset_id = $1
			UNION ALL
			SELECT dd.source_dataset_id, dd.target_dataset_id, dd.task_id, u.depth + 1
			FROM dataset_dependencies dd
			JOIN upstream u ON dd.target_dataset_id = u.source_dataset_id
			WHERE u.depth < 10
		)
		SELECT DISTINCT d.id, d.name, d.version, u.task_id, u.depth
		FROM upstream u
		JOIN datasets d ON d.id = u.source_dataset_id
		ORDER BY u.depth`, datasetID)
	if err != nil {
		return nil, err
	}
	defer upRows.Close()

	var ancestors []map[string]interface{}
	for upRows.Next() {
		var id, taskID, depth int
		var name, version string
		if err := upRows.Scan(&id, &name, &version, &taskID, &depth); err != nil {
			continue
		}
		ancestors = append(ancestors, map[string]interface{}{
			"id": id, "name": name, "version": version, "task_id": taskID, "depth": depth,
		})
	}

	// Находим все потомки (downstream)
	downRows, err := s.DB.Query(`
		WITH RECURSIVE downstream AS (
			SELECT source_dataset_id, target_dataset_id, task_id, 1 AS depth
			FROM dataset_dependencies
			WHERE source_dataset_id = $1
			UNION ALL
			SELECT dd.source_dataset_id, dd.target_dataset_id, dd.task_id, d.depth + 1
			FROM dataset_dependencies dd
			JOIN downstream d ON dd.source_dataset_id = d.target_dataset_id
			WHERE d.depth < 10
		)
		SELECT DISTINCT d.id, d.name, d.version, d2.task_id, d2.depth
		FROM downstream d2
		JOIN datasets d ON d.id = d2.target_dataset_id
		ORDER BY d2.depth`, datasetID)
	if err != nil {
		return nil, err
	}
	defer downRows.Close()

	var descendants []map[string]interface{}
	for downRows.Next() {
		var id, taskID, depth int
		var name, version string
		if err := downRows.Scan(&id, &name, &version, &taskID, &depth); err != nil {
			continue
		}
		descendants = append(descendants, map[string]interface{}{
			"id": id, "name": name, "version": version, "task_id": taskID, "depth": depth,
		})
	}

	// Информация о самом датасете
	var dsName, dsVersion string
	_ = s.DB.QueryRow(`SELECT name, version FROM datasets WHERE id = $1`, datasetID).Scan(&dsName, &dsVersion)

	return map[string]interface{}{
		"dataset": map[string]interface{}{
			"id":      datasetID,
			"name":    dsName,
			"version": dsVersion,
		},
		"ancestors":   ancestors,
		"descendants": descendants,
	}, nil
}

// DeleteTaskDependencies удаляет все зависимости, связанные с задачей.
func (s *DatasetDependencyStore) DeleteTaskDependencies(taskID int) error {
	_, err := s.DB.Exec(`DELETE FROM dataset_dependencies WHERE task_id = $1`, taskID)
	return err
}

// DeleteDatasetDependencies удаляет все зависимости, связанные с датасетом.
func (s *DatasetDependencyStore) DeleteDatasetDependencies(datasetID int) error {
	_, err := s.DB.Exec(`DELETE FROM dataset_dependencies WHERE source_dataset_id = $1 OR target_dataset_id = $1`, datasetID)
	return err
}

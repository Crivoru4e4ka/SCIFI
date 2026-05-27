package services

import (
	"project-MVP/models"
)

// BuildLineageForTask создает зависимости dataset_dependencies для задачи.
// Для каждого output датасета создает связь от каждого input датасета.
func BuildLineageForTask(taskID int, inputIDs, outputIDs []int) error {
	if len(inputIDs) == 0 || len(outputIDs) == 0 {
		return nil
	}
	for _, outID := range outputIDs {
		for _, inID := range inputIDs {
			dep := models.DatasetDependency{
				SourceDatasetID: inID,
				TargetDatasetID: outID,
				TaskID:          taskID,
			}
			if err := DefaultDatasetDependencyStore.CreateDependency(dep); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetDatasetLineage обёртка над DefaultDatasetDependencyStore.
func GetDatasetLineage(datasetID int) (map[string]interface{}, error) {
	return DefaultDatasetDependencyStore.GetLineage(datasetID)
}

// GetDatasetLineageGraph обёртка над DefaultDatasetDependencyStore.
func GetDatasetLineageGraph(datasetID int) (map[string]interface{}, error) {
	return DefaultDatasetDependencyStore.GetFullLineageGraph(datasetID)
}

// DeleteTaskDependencies обёртка над DefaultDatasetDependencyStore.
func DeleteTaskDependencies(taskID int) error {
	return DefaultDatasetDependencyStore.DeleteTaskDependencies(taskID)
}

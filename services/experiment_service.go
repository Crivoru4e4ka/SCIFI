package services

import (
	"project-MVP/models"
)

// GetAllExperimentDatasets обёртка над DefaultExperimentStore.
func GetAllExperimentDatasets() ([]models.ExperimentDataset, error) {
	return DefaultExperimentStore.GetAllExperimentDatasets()
}

// GetDatasetsByTask обёртка над DefaultExperimentStore.
func GetDatasetsByTask(taskID int) ([]models.ExperimentDataset, error) {
	return DefaultExperimentStore.GetDatasetsByTask(taskID)
}

// CreateExperimentDataset обёртка над DefaultExperimentStore.
func CreateExperimentDataset(eds models.ExperimentDataset) (models.ExperimentDataset, error) {
	return DefaultExperimentStore.CreateExperimentDataset(eds)
}

// DeleteExperimentDataset обёртка над DefaultExperimentStore.
func DeleteExperimentDataset(taskID, datasetID int) error {
	return DefaultExperimentStore.DeleteExperimentDataset(taskID, datasetID)
}

// DeleteTaskDatasets обёртка над DefaultExperimentStore.
func DeleteTaskDatasets(taskID int) error {
	return DefaultExperimentStore.DeleteTaskDatasets(taskID)
}

package services

import (
	"project-MVP/models"
)


// GetAllExperimentDatasets обёртка над DefaultExperimentStore.
func GetAllExperimentDatasets() ([]models.ExperimentDataset, error) {
	return DefaultExperimentStore.GetAllExperimentDatasets()
}

// CreateExperimentDataset обёртка над DefaultExperimentStore.
func CreateExperimentDataset(eds models.ExperimentDataset) (models.ExperimentDataset, error) {
	return DefaultExperimentStore.CreateExperimentDataset(eds)
}

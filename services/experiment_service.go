package services

import (
	"project-MVP/db"
	"project-MVP/models"
)

// DefaultExperimentStore — глобальный инстанс ExperimentStore для обратной совместимости.
var DefaultExperimentStore = NewExperimentStore(db.DB)

// GetAllExperimentDatasets обёртка над DefaultExperimentStore.
func GetAllExperimentDatasets() ([]models.ExperimentDataset, error) {
	return DefaultExperimentStore.GetAllExperimentDatasets()
}

// CreateExperimentDataset обёртка над DefaultExperimentStore.
func CreateExperimentDataset(eds models.ExperimentDataset) (models.ExperimentDataset, error) {
	return DefaultExperimentStore.CreateExperimentDataset(eds)
}

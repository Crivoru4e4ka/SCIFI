package services

import (
	"project-MVP/db"
	"project-MVP/models"
)

// DefaultDatasetStore — глобальный инстанс DatasetStore для обратной совместимости.
var DefaultDatasetStore = NewDatasetStore(db.DB)

// GetAllDatasets обёртка над DefaultDatasetStore.
func GetAllDatasets() ([]models.Dataset, error) {
	return DefaultDatasetStore.GetAllDatasets()
}

// CreateDataset обёртка над DefaultDatasetStore.
func CreateDataset(ds models.Dataset) (models.Dataset, error) {
	return DefaultDatasetStore.CreateDataset(ds)
}

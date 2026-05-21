package services

import (
	"project-MVP/models"
)


// GetAllDatasets обёртка над DefaultDatasetStore.
func GetAllDatasets() ([]models.Dataset, error) {
	return DefaultDatasetStore.GetAllDatasets()
}

// CreateDataset обёртка над DefaultDatasetStore.
func CreateDataset(ds models.Dataset) (models.Dataset, error) {
	return DefaultDatasetStore.CreateDataset(ds)
}

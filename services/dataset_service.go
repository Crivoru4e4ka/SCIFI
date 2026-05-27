package services

import (
	"project-MVP/models"
)

// GetAllDatasets обёртка над DefaultDatasetStore.
func GetAllDatasets() ([]models.Dataset, error) {
	return DefaultDatasetStore.GetAllDatasets()
}

// GetDatasetsByProject обёртка над DefaultDatasetStore.
func GetDatasetsByProject(projectID int) ([]models.Dataset, error) {
	return DefaultDatasetStore.GetDatasetsByProject(projectID)
}

// GetDatasetByID обёртка над DefaultDatasetStore.
func GetDatasetByID(id int) (models.Dataset, error) {
	return DefaultDatasetStore.GetDatasetByID(id)
}

// CreateDataset обёртка над DefaultDatasetStore.
func CreateDataset(ds models.Dataset) (models.Dataset, error) {
	return DefaultDatasetStore.CreateDataset(ds)
}

// DeleteDataset обёртка над DefaultDatasetStore.
func DeleteDataset(id int) error {
	return DefaultDatasetStore.DeleteDataset(id)
}

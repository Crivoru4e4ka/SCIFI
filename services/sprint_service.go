package services

import (
	"project-MVP/models"
)


// GetProjectSprints обёртка над DefaultSprintStore.
func GetProjectSprints(projectId int) ([]models.Sprint, error) {
	return DefaultSprintStore.GetProjectSprints(projectId)
}

// CreateSprint обёртка над DefaultSprintStore.
func CreateSprint(s *models.Sprint) error {
	return DefaultSprintStore.CreateSprint(s)
}

// GetSprintByID обёртка над DefaultSprintStore.
func GetSprintByID(id int) (models.Sprint, error) {
	return DefaultSprintStore.GetSprintByID(id)
}

// StartSprint обёртка над DefaultSprintStore.
func StartSprint(id int, s models.Sprint) error {
	return DefaultSprintStore.StartSprint(id, s)
}

// CompleteSprint обёртка над DefaultSprintStore.
func CompleteSprint(sprintID int) error {
	return DefaultSprintStore.CompleteSprint(sprintID)
}

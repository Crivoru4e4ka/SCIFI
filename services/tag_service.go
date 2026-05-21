package services

import (
	"project-MVP/models"
)


// CreateTag обёртка над DefaultTagStore.
func CreateTag(tag models.Tag) (models.Tag, error) {
	return DefaultTagStore.CreateTag(tag)
}

// GetTags обёртка над DefaultTagStore.
func GetTags() ([]models.Tag, error) {
	return DefaultTagStore.GetTags()
}

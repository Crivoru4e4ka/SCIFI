package services

import (
	"project-MVP/db"
	"project-MVP/models"
)

// DefaultTagStore — глобальный инстанс TagStore для обратной совместимости.
var DefaultTagStore = NewTagStore(db.DB)

// CreateTag обёртка над DefaultTagStore.
func CreateTag(tag models.Tag) (models.Tag, error) {
	return DefaultTagStore.CreateTag(tag)
}

// GetTags обёртка над DefaultTagStore.
func GetTags() ([]models.Tag, error) {
	return DefaultTagStore.GetTags()
}

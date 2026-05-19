package services

import (
	"errors"
	"strings"
	"time"

	"project-MVP/db"
	"project-MVP/models"
)

// TagStore инкапсулирует операции с тегами.
type TagStore struct {
	DB db.DBPool
}

// NewTagStore создает новый экземпляр TagStore.
func NewTagStore(database db.DBPool) *TagStore {
	return &TagStore{DB: database}
}

// CreateTag создает новый тег.
func (s *TagStore) CreateTag(tag models.Tag) (models.Tag, error) {
	tag.Name = strings.TrimSpace(tag.Name)
	if tag.Name == "" {
		return models.Tag{}, errors.New("tag name is required")
	}

	createdAt := time.Now()
	query := `INSERT INTO tags (name, created_at) VALUES ($1,$2) RETURNING id`
	var newId int
	if err := s.DB.QueryRow(query, tag.Name, createdAt).Scan(&newId); err != nil {
		return models.Tag{}, err
	}

	tag.ID = newId
	tag.CreatedAt = createdAt
	return tag, nil
}

// GetTags возвращает все теги.
func (s *TagStore) GetTags() ([]models.Tag, error) {
	query := `SELECT id, name, created_at FROM tags`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var t models.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tags, nil
}

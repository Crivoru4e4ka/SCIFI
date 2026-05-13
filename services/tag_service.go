package services

import (
	"errors"
	"project-MVP/db"
	"project-MVP/models"
	"strings"
	"time"
)

func CreateTag(tag models.Tag) (models.Tag, error) {
	tag.Name = strings.TrimSpace(tag.Name)
	if tag.Name == "" {
		return models.Tag{}, errors.New("tag name is required")
	}

	createdAt := time.Now()
	query := `INSERT INTO tags (name, created_at) VALUES ($1,$2) RETURNING id`
	var newId int
	if err := db.DB.QueryRow(query, tag.Name, createdAt).Scan(&newId); err != nil {
		return models.Tag{}, err
	}

	tag.ID = newId
	tag.CreatedAt = createdAt
	return tag, nil
}

func GetTags() ([]models.Tag, error) {
	query := `SELECT id, name, created_at FROM tags`
	rows, err := db.DB.Query(query)
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

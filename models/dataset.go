package models

import "time"

type Dataset struct {
	ID          int         `json:"id"`
	ProjectID   int         `json:"project_id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Version     string      `json:"version"`
	DataURL     string      `json:"data_url"`
	Parameters  interface{} `json:"parameters"`
	CreatedAt   time.Time   `json:"created_at"`
}

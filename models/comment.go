package models

import "time"

type Comment struct {
	Id         int        `json:"id"`
	EntityType string     `json:"entity_type"` // "task" или "project"
	EntityId   int        `json:"entity_id"`
	UserId     int        `json:"user_id"`
	UserName   string     `json:"user_name"`
	Content    string     `json:"content"`
	ParentId   *int       `json:"parent_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
	Replies    []Comment  `json:"replies" swaggertype:"array,object"`
}

package models

import "time"

type Activity struct {
	Id          int       `json:"id"`
	UserId      int       `json:"user_id"`
	UserName    string    `json:"user_name"`
	ProjectId   int       `json:"project_id"`
	ProjectName string    `json:"project_name"`
	EntityType  string    `json:"entity_type"`
	EntityId    int       `json:"entity_id"`
	Action      string    `json:"action"`
	Details     string    `json:"details"`
	CreatedAt   time.Time `json:"created_at"`
}

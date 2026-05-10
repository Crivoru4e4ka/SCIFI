package models

import "time"

type Task struct {
	Id          int        `json:"id"`
	ProjectId   int        `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	AssigneeId  *int       `json:"assignee_id,omitempty"`
	CreatedBy   int        `json:"created_by"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	Team        string     `json:"team,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

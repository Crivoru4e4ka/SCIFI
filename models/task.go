package models

import "time"

type Task struct {
	Id           int       `json:"id"`
	ProjectId    int       `json:"project_id"`
	LocalId      int       `json:"local_id"` // Это наше поле task_num из БД
	SprintId     *int      `json:"sprint_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Status       string    `json:"status"`
	Priority     string    `json:"priority"`
	AssigneeId   *int      `json:"assignee_id"`
	CreatedBy    int       `json:"created_by"`
	DueDate      *string   `json:"due_date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    *string   `json:"updated_at"`
	Type         *string   `json:"type"`
	HypothesisId *int      `json:"hypothesis_id"`
	ResourceId   *int      `json:"resource_id"`
	Conclusion   string    `json:"conclusion"`
	TaskNum      int       `json:"task_num"`
}

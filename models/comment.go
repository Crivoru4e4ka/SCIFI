package models

import "time"

type Comment struct {
	Id        int       `json:"id"`
	TaskId    int       `json:"task_id"`
	UserId    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

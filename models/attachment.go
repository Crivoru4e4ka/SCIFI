package models

import "time"

type Attachment struct {
	Id        int       `json:"id"`
	TaskId    int       `json:"task_id"`
	UserId    int       `json:"user_id"`
	FileName  string    `json:"file_name"`
	FileUrl   string    `json:"file_url"`
	CreatedAt time.Time `json:"created_at"`
}

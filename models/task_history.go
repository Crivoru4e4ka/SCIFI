package models

import "time"

type TaskHistory struct {
	Id        int       `json:"id"`
	TaskId    int       `json:"task_id"`
	UserId    int       `json:"user_id"`
	OldStatus string    `json:"old_status"`
	NewStatus string    `json:"new_status"`
	ChangedAt time.Time `json:"changed_at"`
}

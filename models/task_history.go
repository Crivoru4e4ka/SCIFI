package models

import "time"

type TaskHistory struct {
	Id        int       `json:"id"`
	TaskId    int       `json:"task_id"`
	ChangedBy int       `json:"changed_by"`
	FieldName string    `json:"field_name"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
	ChangedAt time.Time `json:"changed_at"`
}

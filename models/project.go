package models

import "time"

type Project struct {
	Id          int        `json:"id"`
	Name        string     `json:"name"`
	Key         string     `json:"key"`
	Description string     `json:"description"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Status      string     `json:"status"`
	CreatedBy   int        `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
}

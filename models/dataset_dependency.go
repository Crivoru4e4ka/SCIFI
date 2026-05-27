package models

import "time"

type DatasetDependency struct {
	ID              int       `json:"id"`
	SourceDatasetID int       `json:"source_dataset_id"`
	TargetDatasetID int       `json:"target_dataset_id"`
	TaskID          int       `json:"task_id"`
	CreatedAt       time.Time `json:"created_at"`
}

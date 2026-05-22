package models

// ExperimentDataset связывает задачу (task_id) и датасет (dataset_id)
type ExperimentDataset struct {
	TaskID       int    `json:"task_id"`
	DatasetID    int    `json:"dataset_id"`
	RelationType string `json:"relation_type"`
}

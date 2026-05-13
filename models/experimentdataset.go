package models

type ExperimentDataset struct {
	ID           int `json:"id"`
	ExperimentID int `json:"experiment_id"`
	DatasetID    int `json:"dataset_id"`
}

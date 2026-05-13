package handlers

import (
	"encoding/json"
	"net/http"
	"project-MVP/models"
	"project-MVP/services"
)

// GET /experiment-datasets
func GetExperimentDatasets(w http.ResponseWriter, r *http.Request) {
	datasets, err := services.GetAllExperimentDatasets()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(datasets)
}

// POST /experiment-datasets
func CreateExperimentDataset(w http.ResponseWriter, r *http.Request) {
	var eds models.ExperimentDataset
	if err := json.NewDecoder(r.Body).Decode(&eds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	created, err := services.CreateExperimentDataset(eds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

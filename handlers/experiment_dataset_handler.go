package handlers

import (
	"encoding/json"
	"net/http"
	"project-MVP/models"
	"project-MVP/services"
)

// GET /experiment-datasets
func GetExperimentDatasets(w http.ResponseWriter, r *http.Request) {
	_, ok := RequireAuth(w, r)
	if !ok {
		return
	}

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
	_, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	var eds models.ExperimentDataset
	if err := json.NewDecoder(r.Body).Decode(&eds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Experiment-dataset связывает задачу и датасет; проверяем права на задачу
	if eds.TaskID > 0 {
		task, err := services.GetTaskByID(eds.TaskID)
		if err != nil {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		if !RequirePermission(w, r, task.ProjectId, "experiment.create") {
			return
		}
	}

	created, err := services.CreateExperimentDataset(eds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

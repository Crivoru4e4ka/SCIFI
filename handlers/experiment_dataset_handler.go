package handlers

import (
	"encoding/json"
	"net/http"
	"project-MVP/models"
	"project-MVP/services"
	"strconv"

	"github.com/gorilla/mux"
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

// GET /tasks/{id}/datasets
func GetTaskDatasets(w http.ResponseWriter, r *http.Request) {
	_, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil || taskID <= 0 {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	task, err := services.GetTaskByID(taskID)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	if !RequirePermission(w, r, task.ProjectId, "project.view") {
		return
	}

	datasets, err := services.GetDatasetsByTask(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(datasets)
}

// POST /experiment-datasets (общий) и POST /tasks/{id}/datasets (через router)
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

	// Если есть id в URL (/tasks/{id}/datasets), используем его
	vars := mux.Vars(r)
	if taskIDStr, ok := vars["id"]; ok {
		if taskID, err := strconv.Atoi(taskIDStr); err == nil && taskID > 0 {
			eds.TaskID = taskID
		}
	}

	if eds.TaskID <= 0 {
		http.Error(w, "task_id is required", http.StatusBadRequest)
		return
	}

	task, err := services.GetTaskByID(eds.TaskID)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	if !RequirePermission(w, r, task.ProjectId, "task.edit") {
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

// DELETE /tasks/{id}/datasets/{datasetId}
func DeleteTaskDataset(w http.ResponseWriter, r *http.Request) {
	_, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil || taskID <= 0 {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}
	datasetID, err := strconv.Atoi(vars["datasetId"])
	if err != nil || datasetID <= 0 {
		http.Error(w, "invalid dataset id", http.StatusBadRequest)
		return
	}

	task, err := services.GetTaskByID(taskID)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	if !RequirePermission(w, r, task.ProjectId, "task.edit") {
		return
	}

	if err := services.DeleteExperimentDataset(taskID, datasetID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package handlers

import (
	"encoding/json"
	"net/http"
	"project-MVP/models"
	"project-MVP/services"
	"strconv"

	"github.com/gorilla/mux"
)

// GET /datasets
func GetDatasets(w http.ResponseWriter, r *http.Request) {
	_, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	datasets, err := services.GetAllDatasets()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(datasets)
}

// GET /projects/{id}/datasets
func GetProjectDatasets(w http.ResponseWriter, r *http.Request) {
	_, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil || projectID <= 0 {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	if !RequirePermission(w, r, projectID, "project.view") {
		return
	}

	datasets, err := services.GetDatasetsByProject(projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(datasets)
}

// GET /datasets/{id}
func GetDataset(w http.ResponseWriter, r *http.Request) {
	_, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		http.Error(w, "invalid dataset id", http.StatusBadRequest)
		return
	}

	ds, err := services.GetDatasetByID(id)
	if err != nil {
		http.Error(w, "dataset not found", http.StatusNotFound)
		return
	}

	if !RequirePermission(w, r, ds.ProjectID, "project.view") {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ds)
}

// POST /datasets (и POST /projects/{id}/datasets через router)
func CreateDataset(w http.ResponseWriter, r *http.Request) {
	userID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	var ds models.Dataset
	if err := json.NewDecoder(r.Body).Decode(&ds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if ds.ProjectID <= 0 {
		// Попробуем из URL если это /projects/{id}/datasets
		vars := mux.Vars(r)
		if pid, err := strconv.Atoi(vars["id"]); err == nil && pid > 0 {
			ds.ProjectID = pid
		}
	}

	if ds.ProjectID <= 0 {
		http.Error(w, "project_id is required", http.StatusBadRequest)
		return
	}

	if !RequirePermission(w, r, ds.ProjectID, "dataset.upload") {
		return
	}

	ds.CreatedBy = userID

	created, err := services.CreateDataset(ds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

// DELETE /datasets/{id}
func DeleteDataset(w http.ResponseWriter, r *http.Request) {
	_, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		http.Error(w, "invalid dataset id", http.StatusBadRequest)
		return
	}

	ds, err := services.GetDatasetByID(id)
	if err != nil {
		http.Error(w, "dataset not found", http.StatusNotFound)
		return
	}

	if !RequirePermission(w, r, ds.ProjectID, "project.edit") {
		return
	}

	if err := services.DeleteDataset(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

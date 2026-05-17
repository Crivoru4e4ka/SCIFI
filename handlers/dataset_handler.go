package handlers

import (
	"encoding/json"
	"net/http"
	"project-MVP/models"
	"project-MVP/services"
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

// POST /datasets
func CreateDataset(w http.ResponseWriter, r *http.Request) {
	_, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	var ds models.Dataset
	if err := json.NewDecoder(r.Body).Decode(&ds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !RequirePermission(w, r, ds.ProjectID, "dataset.upload") {
		return
	}

	created, err := services.CreateDataset(ds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

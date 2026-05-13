package handlers

import (
	"encoding/json"
	"net/http"
	"project-MVP/models"
	"project-MVP/services"
)

// GET /datasets
func GetDatasets(w http.ResponseWriter, r *http.Request) {
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
	var ds models.Dataset
	if err := json.NewDecoder(r.Body).Decode(&ds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
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

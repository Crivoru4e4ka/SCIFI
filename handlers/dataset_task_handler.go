package handlers

import (
	"encoding/json"
	"net/http"
	"project-MVP/services"
	"strconv"

	"github.com/gorilla/mux"
)

// GET /datasets/{id}/tasks
func GetDatasetTasks(w http.ResponseWriter, r *http.Request) {
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

	tasks, err := services.GetTasksByDataset(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

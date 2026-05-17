package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"project-MVP/models"
	"project-MVP/services"

	"github.com/gorilla/mux"
)

// POST /tasks/{id}/tags
func AddTagToTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil || taskID <= 0 {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	_, ok := RequireAuth(w, r)
	if !ok {
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

	var tag models.Tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, err := services.AddTagToTask(taskID, tag)
	if err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "task or tag not found", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// GET /tasks/{id}/tags
func GetTagsForTask(w http.ResponseWriter, r *http.Request) {
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

	tags, err := services.GetTagsForTask(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

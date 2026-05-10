package handlers

import (
	"encoding/json"
	"net/http"
	"project-MVP/models"
	"project-MVP/services"
	"strconv"

	"github.com/gorilla/mux"
)

// POST /project-members
func CreateProjectMember(w http.ResponseWriter, r *http.Request) {
	var pm models.ProjectMember
	if err := json.NewDecoder(r.Body).Decode(&pm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := services.AddProjectMember(pm)
	if err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "project or user not found", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

// GET /projects/{id}/members
func GetProjectMembersByProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectId, err := strconv.Atoi(vars["id"])
	if err != nil || projectId <= 0 {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	members, err := services.GetProjectMembers(projectId)
	if err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "project not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

package handlers

import (
	"encoding/json"
	"net/http"
	"project-MVP/models"
	"project-MVP/services"
	"strconv"
)

// POST /project-members
func CreateProjectMember(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	var pm models.ProjectMember
	if err := json.NewDecoder(r.Body).Decode(&pm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := services.CheckPermission(userID, pm.ProjectId, "project.manage_members"); err != nil {
		http.Error(w, "insufficient permissions", http.StatusForbidden)
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

	services.LogAudit(userID, pm.ProjectId, "member_added", "project_member", created.Id, "Добавлен участник")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}



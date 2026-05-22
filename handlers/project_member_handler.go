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
	_, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	var pm models.ProjectMember
	if err := json.NewDecoder(r.Body).Decode(&pm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !RequirePermission(w, r, pm.ProjectId, "project.manage_members") {
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

// DELETE /projects/{id}/members/{userID}
func RemoveProjectMemberHandler(w http.ResponseWriter, r *http.Request) {
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
	memberUserID, err := strconv.Atoi(vars["userID"])
	if err != nil || memberUserID <= 0 {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	if !RequirePermission(w, r, projectID, "project.manage_members") {
		return
	}

	if err := services.RemoveProjectMember(projectID, memberUserID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PATCH /projects/{id}/members/{userID}/role
func UpdateProjectMemberRoleHandler(w http.ResponseWriter, r *http.Request) {
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
	memberUserID, err := strconv.Atoi(vars["userID"])
	if err != nil || memberUserID <= 0 {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	if !RequirePermission(w, r, projectID, "project.manage_members") {
		return
	}

	var payload struct {
		RoleID int `json:"role_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем имя роли по role_id
	roles, err := services.GetAllRoles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var roleName string
	for _, role := range roles {
		if role.Id == payload.RoleID {
			roleName = role.Name
			break
		}
	}
	if roleName == "" {
		http.Error(w, "role not found", http.StatusBadRequest)
		return
	}

	roleID := payload.RoleID
	if err := services.UpdateProjectMemberRole(projectID, memberUserID, roleName, &roleID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

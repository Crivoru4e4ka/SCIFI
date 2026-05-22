package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"project-MVP/services"

	"github.com/gorilla/mux"
)

// GetRoles godoc
// @Summary Список всех ролей
// @Description Возвращает список всех доступных системных и проектных ролей
// @Tags rbac
// @Produce json
// @Success 200 {array} models.Role "Список ролей"
// @Failure 500 {string} string "Internal server error"
// @Router /roles [get]
func GetRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := services.GetAllRoles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

// GetPermissions godoc
// @Summary Список всех прав
// @Description Возвращает полный справочник всех прав доступа (permissions), существующих в системе
// @Tags rbac
// @Produce json
// @Success 200 {array} models.Permission "Список прав"
// @Failure 500 {string} string "Internal server error"
// @Router /permissions [get]
func GetPermissions(w http.ResponseWriter, r *http.Request) {
	perms, err := services.GetAllPermissions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(perms)
}

// GetProjectMembersWithRolesHandler godoc
// @Summary Участники проекта с ролями
// @Description Возвращает список всех участников конкретного проекта с их подробными проектными ролями
// @Tags rbac
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {array} models.ProjectMemberWithRole
// @Failure 400 {string} string "invalid project id"
// @Failure 403 {string} string "insufficient permissions"
// @Router /projects/{id}/members [get]
func GetProjectMembersWithRolesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil || projectID <= 0 {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	if !RequirePermission(w, r, projectID, "project.view") {
		return
	}

	members, err := services.GetProjectMembersWithRoles(projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

// AssignProjectRoleHandler godoc
// @Summary Назначить роль участнику
// @Description Изменяет роль пользователя внутри проекта. Требует прав администратора проекта.
// @Tags rbac
// @Accept json
// @Param id path int true "Project ID"
// @Param userID path int true "User ID"
// @Param role body object{role_id=int} true "JSON с ID новой роли"
// @Success 204 "No Content"
// @Failure 400 {string} string "Bad request"
// @Failure 403 {string} string "Forbidden"
// @Router /projects/{id}/members/{userID}/role [post]
func AssignProjectRoleHandler(w http.ResponseWriter, r *http.Request) {
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
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil || userID <= 0 {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	// Проверяем право на управление участниками
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

	if err := services.AssignProjectRole(projectID, userID, payload.RoleID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetMyProjectPermissions godoc
// @Summary Мои права в проекте
// @Description Возвращает список всех кодов прав (permissions), которыми обладает текущий пользователь в данном проекте
// @Tags rbac
// @Param id path int true "Project ID"
// @Produce json
// @Success 200 {object} map[string][]string "Пример: {'permissions': ['task.create', 'task.edit']}"
// @Router /projects/{id}/my-permissions [get]
func GetMyProjectPermissions(w http.ResponseWriter, r *http.Request) {
	userID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil || projectID <= 0 {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	perms, err := services.GetUserPermissions(userID, projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"permissions": perms,
	})
}



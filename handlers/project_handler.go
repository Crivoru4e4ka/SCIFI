package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"project-MVP/models"
	"project-MVP/services"
	"strconv"

	"github.com/gorilla/mux"
)

// CreateProject godoc
// @Summary Создать новый проект
// @Description Создает научный раздел (проект). Создатель автоматически становится руководителем (project_lead).
// @Tags projects
// @Accept json
// @Produce json
// @Param project body models.Project true "Данные проекта"
// @Success 200 {object} models.Project
// @Failure 400 {string} string "Ошибка валидации данных"
// @Failure 401 {string} string "Не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /projects [post]
func CreateProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	var p models.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Принудительно ставим создателя из сессии
	p.CreatedBy = userID

	created, err := services.CreateProject(p)
	if err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "creator not found", http.StatusBadRequest)
			return
		}
		// Валидационные ошибки — 400, остальное — 500
		switch err.Error() {
		case "project name is required", "created_by is required", "team_id is required for team execution type", "project_lead role not found in RBAC system":
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			log.Printf("CreateProject error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

// GetProjects godoc
// @Summary Список всех проектов системы
// @Description Возвращает полный список всех существующих проектов (обычно используется администратором)
// @Tags projects
// @Produce json
// @Success 200 {array} models.Project
// @Failure 500 {string} string "Internal error"
// @Router /projects [get]
func GetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := services.GetProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// GetUserProjects godoc
// @Summary Проекты текущего пользователя
// @Description Возвращает список проектов, к которым у пользователя есть доступ (созданные им, где он участник или открытые проекты)
// @Tags projects
// @Produce json
// @Success 200 {array} models.Project
// @Failure 401 {string} string "Unauthorized"
// @Router /user/projects [get]
func GetUserProjects(w http.ResponseWriter, r *http.Request) {
	userID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	log.Printf("GetUserProjects: Fetching projects for userID: %d\n", userID)

	projects, err := services.GetUserProjects(userID)
	if err != nil {
		log.Printf("GetUserProjects: Error fetching projects: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if projects == nil {
		projects = []models.Project{}
	}

	log.Printf("GetUserProjects: Found %d projects for userID %d\n", len(projects), userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// GetProjectProgress godoc
// @Summary Прогресс выполнения проекта
// @Description Возвращает процент завершенных задач в проекте
// @Tags projects
// @Param id path int true "ID проекта"
// @Produce json
// @Success 200 {object} map[string]interface{} "Пример: {project_id: 1, progress: 75.5}"
// @Router /projects/{id}/progress [get]
func GetProjectProgress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	if !RequirePermission(w, r, id, "project.view") {
		return
	}

	progress, err := services.GetProjectProgress(id)
	if err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "project not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"project_id": id,
		"progress":   progress,
	})
}

// GetProjectTasks godoc
// @Summary Список задач проекта
// @Description Возвращает все задачи, относящиеся к конкретному проекту
// @Tags projects
// @Param id path int true "ID проекта"
// @Produce json
// @Success 200 {array} models.Task
// @Failure 404 {string} string "Project not found"
// @Router /projects/{id}/tasks [get]
func GetProjectTasks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	if !RequirePermission(w, r, id, "project.view") {
		return
	}

	tasks, err := services.GetTasksByProject(id)
	if err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "project not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// GetProjectAssignableUsers godoc
// @Summary Доступные исполнители
// @Description Список пользователей, которых можно назначить на задачи в данном проекте
// @Tags projects
// @Param id path int true "ID проекта"
// @Produce json
// @Success 200 {array} models.User
// @Router /projects/{id}/assignable-users [get]
func GetProjectAssignableUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	if !RequirePermission(w, r, id, "project.view") {
		return
	}

	users, err := services.GetProjectAssignableUsers(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// UpdateProjectHandler godoc
// @Summary Обновить настройки проекта
// @Description Позволяет изменить название, описание, гипотезы и сроки проекта. Требует прав на редактирование.
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "ID проекта"
// @Param project body models.Project true "Объект проекта с новыми данными"
// @Success 204 "No Content"
// @Failure 403 {string} string "Insufficient permissions"
// @Router /projects/{id} [patch]
func UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	userID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	if !RequirePermission(w, r, id, "project.edit") {
		return
	}

	var p models.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := services.UpdateProject(id, p); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	services.LogAudit(userID, id, "project_updated", "project", id, "Изменены настройки проекта")
	w.WriteHeader(http.StatusNoContent)
}

// DeleteProjectHandler godoc
// @Summary Удалить проект
// @Description Полное удаление проекта со всеми задачами. Доступно только руководителю проекта или админу.
// @Tags projects
// @Param id path int true "ID проекта"
// @Success 204 "No Content"
// @Failure 403 {string} string "Insufficient permissions"
// @Router /projects/{id} [delete]
func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	userID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	// Только Lead или Admin
	if !RequirePermission(w, r, id, "project.manage_members") {
		return
	}

	if err := services.DeleteProject(id); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	services.LogActivity(userID, 0, "project", id, "deleted", "Проект удален")
	w.WriteHeader(http.StatusNoContent)
}

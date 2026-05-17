package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project-MVP/models"
	"project-MVP/services"
	"strconv"

	"github.com/gorilla/mux"
)

// POST /tasks — устаревший endpoint, требует авторизации
func CreateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	var t models.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if t.ProjectId <= 0 {
		http.Error(w, "project_id is required", http.StatusBadRequest)
		return
	}

	if !RequirePermission(w, r, t.ProjectId, "task.create") {
		return
	}

	created, err := services.CreateTask(t)
	if err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "project or assignee not found", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	services.LogAudit(userID, t.ProjectId, "task_created", "task", created.Id, "Создана задача")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

// POST /projects/{id}/tasks
func CreateTaskInProject(w http.ResponseWriter, r *http.Request) {
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

	if !RequirePermission(w, r, projectID, "task.create") {
		return
	}

	var t models.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t.ProjectId = projectID

	created, err := services.CreateTask(t)
	if err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "project or assignee not found", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	services.LogAudit(userID, projectID, "task_created", "task", created.Id, "Создана задача")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

// PATCH /tasks/{id}/status
func UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	currentUserID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	var data struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем задачу для проверки project_id
	task, err := services.GetTaskByID(id)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	// Проверяем права на изменение статуса
	if !RequirePermission(w, r, task.ProjectId, "task.change_status") {
		return
	}

	if err := services.UpdateTaskStatus(id, currentUserID, data.Status); err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	services.LogAudit(currentUserID, task.ProjectId, "task_status_changed", "task", id,
		fmt.Sprintf("Статус изменен на %s", data.Status))
	w.WriteHeader(http.StatusNoContent)
}

// PATCH /tasks/{id}/sprint
func UpdateTaskSprintHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	_, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	// Получаем задачу для проверки project_id
	task, err := services.GetTaskByID(id)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	if !RequirePermission(w, r, task.ProjectId, "task.edit") {
		return
	}

	var data struct {
		SprintId *int `json:"sprint_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := services.UpdateTaskSprint(id, data.SprintId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetAllUserTasksHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	tasks, err := services.GetAllUserTasks(userID)
	if err != nil {
		http.Error(w, "Ошибка при получении всех задач: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// POST /projects/{id}/hypotheses
func CreateHypothesisHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	projectID, _ := strconv.Atoi(vars["id"])

	if !RequirePermission(w, r, projectID, "hypothesis.edit") {
		return
	}

	var h models.Hypothesis
	if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	h.ProjectID = projectID

	newHypo, err := services.CreateHypothesis(h)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newHypo)
}

// GET /projects/{id}/hypotheses
func GetProjectHypothesesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, _ := strconv.Atoi(vars["id"])

	if !RequirePermission(w, r, projectID, "project.view") {
		return
	}

	list, err := services.GetProjectHypotheses(projectID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

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

// CreateTask godoc
// @Summary Создать задачу (базовый метод)
// @Description Создает задачу. Рекомендуется использовать CreateTaskInProject для привязки к разделу.
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Объект задачи"
// @Success 200 {object} models.Task
// @Failure 400 {string} string "Bad request"
// @Router /tasks [post]
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

// CreateTaskInProject godoc
// @Summary Создать задачу в проекте
// @Description Создает новую научную задачу, привязанную к конкретному разделу (проекту)
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param task body models.Task true "Данные задачи"
// @Success 200 {object} models.Task
// @Failure 400 {string} string "invalid project id"
// @Failure 403 {string} string "forbidden"
// @Router /projects/{id}/tasks [post]
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

// UpdateTaskStatus godoc
// @Summary Изменить статус задачи
// @Description Позволяет перевести задачу в другой статус (например, из 'В работе' в 'На проверке')
// @Tags tasks
// @Accept json
// @Param id path int true "Task ID"
// @Param status body object{status=string} true "JSON со статусом"
// @Success 204 "No Content"
// @Failure 403 {string} string "forbidden"
// @Failure 404 {string} string "task not found"
// @Router /tasks/{id}/status [patch]
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

// UpdateTaskSprintHandler godoc
// @Summary Переместить задачу в спринт или бэклог
// @Description Привязывает задачу к указанному спринту. Если передать null, задача вернется в бэклог проекта.
// @Tags tasks
// @Accept json
// @Param id path int true "Task ID"
// @Param body body object{sprint_id=int} true "ID Спринта"
// @Success 204 "No Content"
// @Router /tasks/{id}/sprint [patch]
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

// GetAllUserTasksHandler godoc
// @Summary Получить все задачи текущего пользователя
// @Description Список всех задач из всех проектов, где текущий пользователь является участником
// @Tags tasks
// @Produce json
// @Success 200 {array} models.Task
// @Router /user/tasks [get]
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

// CreateHypothesisHandler godoc
// @Summary Создать гипотезу в проекте
// @Description Добавляет новую научную гипотезу, которую необходимо проверить в рамках исследования
// @Tags hypotheses
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param hypothesis body models.Hypothesis true "Данные гипотезы"
// @Success 200 {object} models.Hypothesis
// @Router /projects/{id}/hypotheses [post]
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

// GetProjectHypothesesHandler godoc
// @Summary Список гипотез проекта
// @Tags hypotheses
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {array} models.Hypothesis
// @Router /projects/{id}/hypotheses [get]
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

// UpdateTaskHandler godoc
// @Summary Редактировать данные задачи
// @Description Позволяет изменить описание, приоритет, тип и другие параметры задачи. Требует прав на редактирование.
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body models.Task true "Объект задачи с новыми данными"
// @Success 204 "No Content"
// @Failure 404 {string} string "Task not found"
// @Router /tasks/{id} [patch]
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	userID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	var t models.Task
	json.NewDecoder(r.Body).Decode(&t)

	oldTask, err := services.GetTaskByID(id)
	if err != nil {
		http.Error(w, "Task not found", 404)
		return
	}

	if !RequirePermission(w, r, oldTask.ProjectId, "task.edit") {
		return
	}

	if err := services.UpdateTask(id, t); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	services.LogAudit(userID, oldTask.ProjectId, "task_updated", "task", id, "Обновлена информация о задаче")
	w.WriteHeader(http.StatusNoContent)
}

// DeleteTaskHandler godoc
// @Summary Удалить задачу
// @Description Полное удаление задачи. Доступно только автору или руководителю проекта.
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 204 "No Content"
// @Failure 403 {string} string "Forbidden"
// @Router /tasks/{id} [delete]
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	userID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	task, err := services.GetTaskByID(id)
	if err != nil {
		http.Error(w, "Task not found", 404)
		return
	}

	// Удалять может автор или Lead
	isLead := services.CheckPermission(userID, task.ProjectId, "project.manage_members") == nil
	if task.CreatedBy != userID && !isLead {
		http.Error(w, "Forbidden", 403)
		return
	}

	if err := services.DeleteTask(id); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	services.LogAudit(userID, task.ProjectId, "task_deleted", "task", id, "Задача удалена")
	w.WriteHeader(http.StatusNoContent)
}

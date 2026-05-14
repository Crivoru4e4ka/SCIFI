package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"project-MVP/models"
	"project-MVP/services"

	"github.com/gorilla/mux"
)

// POST /tasks
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var t models.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

// POST /projects/{id}/tasks
func CreateTaskInProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil || projectID <= 0 {
		http.Error(w, "invalid project id", http.StatusBadRequest)
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

	// ДОБАВЛЕНО: Получаем ID пользователя из сессии
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	currentUserID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	var data struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ИСПРАВЛЕНО: Теперь передаем три аргумента
	if err := services.UpdateTaskStatus(id, currentUserID, data.Status); err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateTaskSprintHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

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
	// 1. Получаем сессию пользователя (как в ваших прошлых хендлерах)
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}
	userID, _ := strconv.Atoi(cookie.Value)

	// 2. Вызываем сервис
	tasks, err := services.GetAllUserTasks(userID)
	if err != nil {
		http.Error(w, "Ошибка при получении всех задач: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Отправляем JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

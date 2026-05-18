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

// GetProjectSprintsHandler godoc
// @Summary Список спринтов проекта
// @Description Возвращает все спринты (запланированные, активные и завершенные) для конкретного проекта
// @Tags sprints
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {array} models.Sprint "Список спринтов"
// @Failure 400 {string} string "Некорректный ID"
// @Failure 403 {string} string "Нет прав доступа"
// @Router /projects/{id}/sprints [get]
func GetProjectSprintsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Некорректный ID проекта", http.StatusBadRequest)
		return
	}

	if !RequirePermission(w, r, projectID, "project.view") {
		return
	}

	sprints, err := services.GetProjectSprints(projectID)
	if err != nil {
		http.Error(w, "Ошибка при получении спринтов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sprints); err != nil {
		http.Error(w, "Ошибка кодирования JSON", http.StatusInternalServerError)
	}
}

// CreateSprintHandler godoc
// @Summary Создать новый спринт
// @Description Создает новый спринт в статусе 'planned' внутри указанного проекта
// @Tags sprints
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param sprint body models.Sprint true "Данные спринта (только имя)"
// @Success 201 {object} models.Sprint "Созданный спринт"
// @Failure 403 {string} string "Нет прав на управление спринтами"
// @Router /projects/{id}/sprints [post]
func CreateSprintHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Некорректный ID проекта", http.StatusBadRequest)
		return
	}

	if !RequirePermission(w, r, projectID, "sprint.manage") {
		return
	}

	var s models.Sprint
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Ошибка в формате данных", http.StatusBadRequest)
		return
	}

	s.ProjectID = projectID
	if s.Status == "" {
		s.Status = "planned"
	}

	if err := services.CreateSprint(&s); err != nil {
		http.Error(w, "Не удалось создать спринт: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

// StartSprintHandler godoc
// @Summary Запустить спринт
// @Description Активирует спринт, устанавливая даты начала, окончания и цель исследования на данный период
// @Tags sprints
// @Accept json
// @Param id path int true "Sprint ID"
// @Param data body models.Sprint true "Данные запуска (Name, StartDate, EndDate, Goal)"
// @Success 200 {string} string "OK"
// @Failure 404 {string} string "Sprint not found"
// @Router /sprints/{id}/start [patch]
func StartSprintHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	_, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	// Получаем спринт для проверки project_id
	sprint, err := services.GetSprintByID(id)
	if err != nil {
		http.Error(w, "sprint not found", http.StatusNotFound)
		return
	}

	if !RequirePermission(w, r, sprint.ProjectID, "sprint.manage") {
		return
	}

	var s models.Sprint
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := services.StartSprint(id, s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// CompleteSprintHandler godoc
// @Summary Завершить спринт
// @Description Переводит спринт в статус 'completed'. Задачи, которые не были выполнены, автоматически возвращаются в бэклог.
// @Tags sprints
// @Param id path int true "Sprint ID"
// @Success 204 "No Content"
// @Failure 404 {string} string "Sprint not found"
// @Router /sprints/{id}/complete [patch]
func CompleteSprintHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	_, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	// Получаем спринт для проверки project_id
	sprint, err := services.GetSprintByID(id)
	if err != nil {
		http.Error(w, "sprint not found", http.StatusNotFound)
		return
	}

	if !RequirePermission(w, r, sprint.ProjectID, "sprint.manage") {
		return
	}

	if err := services.CompleteSprint(id); err != nil {
		log.Printf("Ошибка при завершении спринта: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

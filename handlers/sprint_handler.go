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
	w.WriteHeader(http.StatusNoContent)
}

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

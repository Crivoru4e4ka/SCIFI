package handlers

import (
	"encoding/json"
	"net/http"
	"project-MVP/models"
	"project-MVP/services"
	"strconv"

	"github.com/gorilla/mux"
)

func GetProjectSprintsHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Получаем ID проекта из параметров URL
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Некорректный ID проекта", http.StatusBadRequest)
		return
	}

	// 2. Вызываем сервис (БЕЗ передачи sql.DB, так как сервис сам знает о базе)
	sprints, err := services.GetProjectSprints(projectID)
	if err != nil {
		// Выводим ошибку в консоль сервера для отладки
		http.Error(w, "Ошибка при получении спринтов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Устанавливаем заголовок и отправляем JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sprints); err != nil {
		http.Error(w, "Ошибка кодирования JSON", http.StatusInternalServerError)
	}
}

func CreateSprintHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Получаем ID проекта из URL
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Некорректный ID проекта", http.StatusBadRequest)
		return
	}

	// 2. Читаем JSON, который прислал фронтенд
	var s models.Sprint
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Ошибка в формате данных", http.StatusBadRequest)
		return
	}

	// 3. Назначаем ID проекта и статус по умолчанию
	s.ProjectID = projectID
	if s.Status == "" {
		s.Status = "planned"
	}

	// 4. Сохраняем в базу через сервис
	if err := services.CreateSprint(&s); err != nil {
		http.Error(w, "Не удалось создать спринт: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Возвращаем созданный спринт (теперь уже с ID) обратно фронтенду
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func StartSprintHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

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

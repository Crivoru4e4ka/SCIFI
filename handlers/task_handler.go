package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project-MVP/db"
	"project-MVP/models"
	"project-MVP/services"
	"strconv"

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

func GetProjectReportData(projectID int) (string, error) {
	var projectName string
	var projectDesc string
	// Получаем данные проекта
	err := db.DB.QueryRow("SELECT name, description FROM projects WHERE id = $1", projectID).Scan(&projectName, &projectDesc)
	if err != nil {
		return "", err
	}

	// Получаем все завершенные задачи
	query := `SELECT task_num, title, COALESCE(description, ''), COALESCE(conclusion, '') 
	          FROM tasks WHERE project_id = $1 AND status IN ('done', 'ГОТОВО') ORDER BY task_num`
	rows, err := db.DB.Query(query, projectID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// Формируем текст по структуре ГОСТ 7.32
	report := "ОТЧЕТ О НАУЧНО-ИССЛЕДОВАТЕЛЬСКОЙ РАБОТЕ\n"
	report += "Тема: " + projectName + "\n\n"
	report += "1. ВВЕДЕНИЕ\n"
	report += projectDesc + "\n\n"
	report += "2. ОСНОВНАЯ ЧАСТЬ (РЕЗУЛЬТАТЫ ЭТАПОВ)\n"

	for rows.Next() {
		var num int
		var title, desc, conc string
		rows.Scan(&num, &title, &desc, &conc)
		report += fmt.Sprintf("\nЭтап %d: %s\n", num, title)
		report += "Описание работ: " + desc + "\n"
		if conc != "" {
			report += "Научный вывод: " + conc + "\n"
		}
	}

	report += "\n\n3. ЗАКЛЮЧЕНИЕ\n"
	report += "Задачи этапа НИР выполнены в полном объеме."

	return report, nil
}

// POST /projects/{id}/hypotheses
func CreateHypothesisHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, _ := strconv.Atoi(vars["id"])

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

	list, err := services.GetProjectHypotheses(projectID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

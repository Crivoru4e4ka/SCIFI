package handlers

import (
	"encoding/json"
	"net/http"
	"project-MVP/models"
	"project-MVP/services"
	"strconv"

	"github.com/gorilla/mux"
)

// GetUserTeamsHandler — получает все команды, в которых состоит пользователь
func GetUserTeamsHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Достаем ID пользователя из URL
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Некорректный ID пользователя", http.StatusBadRequest)
		return
	}

	// 2. Вызываем сервис для получения данных из БД
	teams, err := services.GetUserTeams(userID)
	if err != nil {
		http.Error(w, "Ошибка при получении команд: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Отправляем результат в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

// CreateTeamHandler — создает новую команду и добавляет создателя в её участники
func CreateTeamHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Получаем ID текущего пользователя из куки (сессии)
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "Нужна авторизация", http.StatusUnauthorized)
		return
	}
	currentUserID, _ := strconv.Atoi(cookie.Value)

	// 2. Читаем данные из тела запроса
	var t models.Team
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Ошибка в формате данных", http.StatusBadRequest)
		return
	}
	t.CreatedBy = currentUserID

	// 3. Вызываем сервис создания (нужно будет добавить в services)
	// Для диплома важно, чтобы создатель сразу стал участником команды в team_members
	newTeam, err := services.CreateTeam(t)
	if err != nil {
		http.Error(w, "Не удалось создать команду", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTeam)
}

func GetTeamMembersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID, _ := strconv.Atoi(vars["id"])

	// Вызываем сервис
	members, err := services.GetTeamMembers(teamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

func RemoveTeamMemberHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID, _ := strconv.Atoi(vars["id"])
	userID, _ := strconv.Atoi(vars["userID"])

	if err := services.RemoveMemberFromTeam(teamID, userID); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

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

// GetUserTeamsHandler — получает все команды, в которых состоит пользователь
func GetUserTeamsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Некорректный ID пользователя", http.StatusBadRequest)
		return
	}

	// Позволяем смотреть свои команды или требуем админа
	currentUserID, ok := GetUserID(r)
	if !ok {
		http.Error(w, "Нужна авторизация", http.StatusUnauthorized)
		return
	}
	if currentUserID != userID && !services.IsAdmin(currentUserID) {
		http.Error(w, "insufficient permissions", http.StatusForbidden)
		return
	}

	teams, err := services.GetUserTeams(userID)
	if err != nil {
		http.Error(w, "Ошибка при получении команд: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

// CreateTeamHandler — создает новую команду и добавляет создателя в её участники
func CreateTeamHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Получаем ID текущего пользователя из контекста
	currentUserID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	// 2. Читаем данные из тела запроса
	var t models.Team
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Ошибка в формате данных", http.StatusBadRequest)
		return
	}
	t.CreatedBy = currentUserID

	// 3. Вызываем сервис создания
	newTeam, err := services.CreateTeam(t)
	if err != nil {
		log.Printf("CreateTeam error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTeam)
}

func GetTeamMembersHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	teamID, _ := strconv.Atoi(vars["id"])

	members, err := services.GetTeamMembers(teamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

func RemoveTeamMemberHandler(w http.ResponseWriter, r *http.Request) {
	currentUserID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	teamID, _ := strconv.Atoi(vars["id"])
	userID, _ := strconv.Atoi(vars["userID"])

	// Проверяем, что текущий пользователь — создатель команды или админ
	team, err := services.GetTeamByID(teamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if team.CreatedBy != currentUserID && !services.IsAdmin(currentUserID) {
		http.Error(w, "insufficient permissions", http.StatusForbidden)
		return
	}

	if err := services.RemoveMemberFromTeam(teamID, userID); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// PATCH /teams/{id}
func UpdateTeamHandler(w http.ResponseWriter, r *http.Request) {
	currentUserID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	teamID, _ := strconv.Atoi(vars["id"])

	team, err := services.GetTeamByID(teamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if team.CreatedBy != currentUserID && !services.IsAdmin(currentUserID) {
		http.Error(w, "insufficient permissions", http.StatusForbidden)
		return
	}

	var req models.UpdateTeamRequest
	json.NewDecoder(r.Body).Decode(&req)

	if err := services.UpdateTeam(teamID, req.Name, req.Description); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// POST /teams/{id}/members
func AddTeamMemberHandler(w http.ResponseWriter, r *http.Request) {
	currentUserID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	teamID, _ := strconv.Atoi(vars["id"])

	team, err := services.GetTeamByID(teamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if team.CreatedBy != currentUserID && !services.IsAdmin(currentUserID) {
		http.Error(w, "insufficient permissions", http.StatusForbidden)
		return
	}

	var req models.AddMemberRequest
	json.NewDecoder(r.Body).Decode(&req)

	member, err := services.AddMemberToTeam(teamID, req.Email)
	if err != nil {
		http.Error(w, "Пользователь не найден", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

// PATCH /teams/{id}/members/{userID}/role
func UpdateMemberRoleHandler(w http.ResponseWriter, r *http.Request) {
	currentUserID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	teamID, _ := strconv.Atoi(vars["id"])
	userID, _ := strconv.Atoi(vars["userID"])

	team, err := services.GetTeamByID(teamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if team.CreatedBy != currentUserID && !services.IsAdmin(currentUserID) {
		http.Error(w, "insufficient permissions", http.StatusForbidden)
		return
	}

	var req models.UpdateRoleRequest
	json.NewDecoder(r.Body).Decode(&req)

	if err := services.UpdateMemberRole(teamID, userID, req.Role); err != nil {
		if err.Error() == "invalid team member role: must be a valid project role" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DELETE /teams/{id}
func DeleteTeamHandler(w http.ResponseWriter, r *http.Request) {
	currentUserID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	teamID, _ := strconv.Atoi(vars["id"])

	team, err := services.GetTeamByID(teamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if team.CreatedBy != currentUserID && !services.IsAdmin(currentUserID) {
		http.Error(w, "insufficient permissions", http.StatusForbidden)
		return
	}

	if err := services.DeleteTeam(teamID); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

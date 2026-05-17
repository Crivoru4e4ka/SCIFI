package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"project-MVP/db"
	"project-MVP/models"
	"project-MVP/services"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	var req models.Comment
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Определяем ID проекта для логов
	projectID := req.EntityId
	if req.EntityType == "task" {
		err := db.DB.QueryRow("SELECT project_id FROM tasks WHERE id = $1", req.EntityId).Scan(&projectID)
		if err != nil {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
	}

	// Проверка прав (состоит ли пользователь в проекте)
	if err := services.CheckPermission(userID, projectID, "project.view"); err != nil {
		http.Error(w, "no access to this project", http.StatusForbidden)
		return
	}

	req.UserId = userID
	id, err := services.AddComment(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	details := "Добавлен новый комментарий"
	if req.ParentId != nil {
		details = "Добавлен ответ на комментарий"
	}
	if req.EntityType == "project" {
		details += " к разделу"
	} else {
		details += " к задаче"
	}

	log.Printf("[DEBUG] Logging comment for Project: %d, Action: comment_created", projectID)

	services.LogAudit(userID, projectID, "comment_created", "comment", id, details)
	// Записываем в таблицу активности
	services.LogActivity(userID, projectID, "comment", id, "created", details)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eType := vars["type"] // "task" или "project"
	eId, _ := strconv.Atoi(vars["id"])

	list, err := services.GetCommentsByEntity(eType, eId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	commentID, _ := strconv.Atoi(vars["id"])

	comment, err := services.GetCommentRaw(commentID)
	if err != nil {
		http.Error(w, "comment not found", http.StatusNotFound)
		return
	}

	// Находим проект для проверки прав модератора
	projectID := comment.EntityId
	if comment.EntityType == "task" {
		db.DB.QueryRow("SELECT project_id FROM tasks WHERE id = $1", comment.EntityId).Scan(&projectID)
	}

	// ПРАВА: Удалить может автор ИЛИ тот, у кого есть право project.manage_members (Lead)
	isAuthor := comment.UserId == userID
	isLead := services.CheckPermission(userID, projectID, "project.manage_members") == nil

	if !isAuthor && !isLead {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := services.SoftDeleteComment(commentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	services.LogActivity(userID, projectID, "comment", commentID, "deleted", "Пользователь удалил свой комментарий")
	services.LogAudit(userID, projectID, "comment_deleted", "comment", commentID, "Пользователь удалил свой комментарий")

	w.WriteHeader(http.StatusNoContent)
}

// PATCH /comments/{id}
func UpdateCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	userID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	comment, _ := services.GetCommentRaw(id)
	if comment.UserId != userID {
		http.Error(w, "Not your comment", 403)
		return
	}

	var data struct {
		Content string `json:"content"`
	}
	json.NewDecoder(r.Body).Decode(&data)

	services.UpdateComment(id, data.Content)
	w.WriteHeader(200)
}

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

// CreateCommentHandler godoc
// @Summary Добавить комментарий
// @Description Создает новый основной комментарий или ответ на существующий (через parent_id). Доступно только участникам проекта.
// @Tags comments
// @Accept json
// @Produce json
// @Param comment body models.Comment true "Объект комментария"
// @Success 201 {object} map[string]int "ID созданного комментария"
// @Failure 400 {string} string "invalid request"
// @Failure 403 {string} string "no access to this project"
// @Failure 500 {string} string "internal error"
// @Router /comments [post]
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

// GetCommentsHandler godoc
// @Summary Получить дерево комментариев
// @Description Возвращает иерархический список комментариев для конкретной задачи или всего проекта
// @Tags comments
// @Param type path string true "Тип сущности (task или project)"
// @Param id path int true "ID сущности"
// @Produce json
// @Success 200 {array} models.Comment "Дерево комментариев"
// @Failure 500 {string} string "internal error"
// @Router /comments/{type}/{id} [get]
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

// DeleteCommentHandler godoc
// @Summary Удалить комментарий
// @Description Выполняет мягкое удаление (soft delete). Удалить может автор комментария или руководитель проекта (project_lead).
// @Tags comments
// @Param id path int true "ID комментария"
// @Success 204 "No Content"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 404 {string} string "comment not found"
// @Router /comments/{id} [delete]
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

// UpdateCommentHandler godoc
// @Summary Редактировать комментарий
// @Description Позволяет автору изменить текст своего комментария
// @Tags comments
// @Accept json
// @Param id path int true "ID комментария"
// @Param content body object{content=string} true "Новое содержимое"
// @Success 200 "OK"
// @Failure 403 {string} string "Not your comment"
// @Router /comments/{id} [patch]
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

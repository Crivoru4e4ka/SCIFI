package handlers

import (
	"encoding/json"
	"net/http"
	"project-MVP/models"
	"project-MVP/services"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type createCommentRequest struct {
	UserId  int    `json:"user_id"`
	Content string `json:"content"`
}

// POST /tasks/{id}/comments
func CreateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["id"])
	if err != nil || taskId <= 0 {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	var req createCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comment := models.Comment{
		TaskId:    taskId,
		UserId:    req.UserId,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	created, err := services.CreateComment(comment)
	if err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "task or user not found", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

// GET /tasks/{id}/comments
func GetCommentsByTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["id"])
	if err != nil || taskId <= 0 {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	comments, err := services.GetCommentsByTask(taskId)
	if err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

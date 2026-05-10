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

// POST /projects
func CreateProject(w http.ResponseWriter, r *http.Request) {
	var p models.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := services.CreateProject(p)
	if err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "creator not found", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

// GET /projects
func GetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := services.GetProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// GET /user/projects - получить проекты текущего пользователя
func GetUserProjects(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println("GetUserProjects: No session cookie found")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		log.Printf("GetUserProjects: Invalid session value: %s\n", cookie.Value)
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	log.Printf("GetUserProjects: Fetching projects for userID: %d\n", userID)

	projects, err := services.GetUserProjects(userID)
	if err != nil {
		log.Printf("GetUserProjects: Error fetching projects: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if projects == nil {
		projects = []models.Project{}
	}

	log.Printf("GetUserProjects: Found %d projects for userID %d\n", len(projects), userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// GET /projects/{id}/progress
func GetProjectProgress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	progress, err := services.GetProjectProgress(id)
	if err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "project not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"project_id": id,
		"progress":   progress,
	})
}

// GET /projects/{id}/tasks
func GetProjectTasks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	tasks, err := services.GetTasksByProject(id)
	if err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "project not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

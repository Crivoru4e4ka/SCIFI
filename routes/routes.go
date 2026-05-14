package routes

import (
	"net/http"
	"project-MVP/handlers"

	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// --- Auth ---
	r.HandleFunc("/auth/register", handlers.Register).Methods("POST")
	r.HandleFunc("/auth/login", handlers.Login).Methods("POST")

	// --- Users ---
	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/me", handlers.GetCurrentUser).Methods("GET")

	// --- Projects ---
	r.HandleFunc("/projects", handlers.CreateProject).Methods("POST")
	r.HandleFunc("/projects", handlers.GetProjects).Methods("GET")
	r.HandleFunc("/user/projects", handlers.GetUserProjects).Methods("GET")
	r.HandleFunc("/projects/{id}/tasks", handlers.CreateTaskInProject).Methods("POST")
	r.HandleFunc("/projects/{id}/tasks", handlers.GetProjectTasks).Methods("GET")
	r.HandleFunc("/projects/{id}/progress", handlers.GetProjectProgress).Methods("GET")

	// --- Tasks ---
	r.HandleFunc("/tasks", handlers.CreateTask).Methods("POST")
	r.HandleFunc("/tasks/{id}/status", handlers.UpdateTaskStatus).Methods("PATCH")
	r.HandleFunc("/tasks/{id}/comments", handlers.CreateComment).Methods("POST")
	r.HandleFunc("/tasks/{id}/comments", handlers.GetCommentsByTask).Methods("GET")
	r.HandleFunc("/tasks/{id}/sprint", handlers.UpdateTaskSprintHandler).Methods("PATCH")
	r.HandleFunc("/user/tasks", handlers.GetAllUserTasksHandler).Methods("GET")

	// --- Sprints ---
	r.HandleFunc("/projects/{id}/sprints", handlers.GetProjectSprintsHandler).Methods("GET")
	r.HandleFunc("/projects/{id}/sprints", handlers.CreateSprintHandler).Methods("POST")
	r.HandleFunc("/sprints/{id}/start", handlers.StartSprintHandler).Methods("PATCH")
	r.HandleFunc("/sprints/{id}/complete", handlers.CompleteSprintHandler).Methods("PATCH")

	// --- Attachments ---
	r.HandleFunc("/projects/{id}/attachments", handlers.GetProjectAttachmentsHandler).Methods("GET")
	r.HandleFunc("/tasks/{id}/attachments", handlers.UploadFileHandler).Methods("POST")

	// --- Project Members ---
	r.HandleFunc("/project-members", handlers.CreateProjectMember).Methods("POST")
	r.HandleFunc("/projects/{id}/members", handlers.GetProjectMembersByProject).Methods("GET")

	// --- Команды ---
	r.HandleFunc("/user/{id}/teams", handlers.GetUserTeamsHandler).Methods("GET")
	r.HandleFunc("/teams", handlers.CreateTeamHandler).Methods("POST")
	r.HandleFunc("/teams/{id}", handlers.UpdateTeamHandler).Methods("PATCH")
	r.HandleFunc("/teams/{id}", handlers.DeleteTeamHandler).Methods("DELETE")
	r.HandleFunc("/teams/{id}/members", handlers.GetTeamMembersHandler).Methods("GET")
	r.HandleFunc("/teams/{id}/members", handlers.AddTeamMemberHandler).Methods("POST")
	r.HandleFunc("/teams/{id}/members/{userID}", handlers.RemoveTeamMemberHandler).Methods("DELETE")
	r.HandleFunc("/teams/{id}/members/{userID}/role", handlers.UpdateMemberRoleHandler).Methods("PATCH")

	// --- UI ---
	r.HandleFunc("/", handlers.ServeIndex).Methods("GET")
	r.HandleFunc("/login", handlers.ServeLogin).Methods("GET")
	r.HandleFunc("/logout", handlers.Logout).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web"))))

	// --- Отчеты ---
	r.HandleFunc("/projects/{id}/report", handlers.GenerateGostReport).Methods("GET")

	return r
}

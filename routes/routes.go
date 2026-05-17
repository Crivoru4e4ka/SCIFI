package routes

import (
	"net/http"
	"project-MVP/handlers"
	"project-MVP/middleware"

	"github.com/gorilla/mux"
)

// NoCacheMiddleware добавляет заголовки, запрещающие браузеру кэшировать страницы.
// Это критически важно для безопасности личного кабинета.
func NoCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}

func InitRoutes() *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.AuthMiddleware)
	r.Use(NoCacheMiddleware)

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
	r.HandleFunc("/projects/{id}/hypotheses", handlers.GetProjectHypothesesHandler).Methods("GET")
	r.HandleFunc("/projects/{id}/hypotheses", handlers.CreateHypothesisHandler).Methods("POST")
	r.HandleFunc("/projects/{id}", handlers.UpdateProjectHandler).Methods("PATCH")
	r.HandleFunc("/projects/{id}", handlers.DeleteProjectHandler).Methods("DELETE")

	// --- Tasks ---
	r.HandleFunc("/tasks", handlers.CreateTask).Methods("POST")
	r.HandleFunc("/tasks/{id}/status", handlers.UpdateTaskStatus).Methods("PATCH")
	r.HandleFunc("/tasks/{id}/sprint", handlers.UpdateTaskSprintHandler).Methods("PATCH")
	r.HandleFunc("/user/tasks", handlers.GetAllUserTasksHandler).Methods("GET")
	r.HandleFunc("/tasks/{id}", handlers.UpdateTaskHandler).Methods("PATCH")
	r.HandleFunc("/tasks/{id}", handlers.DeleteTaskHandler).Methods("DELETE")

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
	r.HandleFunc("/projects/{id}/members", handlers.GetProjectMembersWithRolesHandler).Methods("GET")
	r.HandleFunc("/projects/{id}/members/{userID}", handlers.RemoveProjectMemberHandler).Methods("DELETE")
	r.HandleFunc("/projects/{id}/members/{userID}/role", handlers.UpdateProjectMemberRoleHandler).Methods("PATCH")
	r.HandleFunc("/projects/{id}/assignable-users", handlers.GetProjectAssignableUsers).Methods("GET")

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

	// --- Активности ---
	r.HandleFunc("/activities", handlers.GetActivitiesHandler).Methods("GET")

	// --- Комментарии ---
	r.HandleFunc("/comments", handlers.CreateCommentHandler).Methods("POST")
	r.HandleFunc("/comments/{type}/{id}", handlers.GetCommentsHandler).Methods("GET")
	r.HandleFunc("/comments/{id}", handlers.DeleteCommentHandler).Methods("DELETE")
	r.HandleFunc("/comments/{id}", handlers.UpdateCommentHandler).Methods("PATCH")

	// --- Grants ---
	r.HandleFunc("/grants", handlers.CreateGrant).Methods("POST")
	r.HandleFunc("/grants", handlers.ListGrants).Methods("GET")
	r.HandleFunc("/grants/{id}", handlers.GetGrant).Methods("GET")
	r.HandleFunc("/grants/{id}", handlers.UpdateGrant).Methods("PATCH")
	r.HandleFunc("/grants/{id}", handlers.DeleteGrant).Methods("DELETE")
	r.HandleFunc("/grants/{id}/projects", handlers.GetGrantProjects).Methods("GET")
	r.HandleFunc("/grants/{id}/projects", handlers.AddProjectToGrant).Methods("POST")
	r.HandleFunc("/grants/{id}/budget", handlers.GetGrantBudget).Methods("GET")
	r.HandleFunc("/grants/{grantId}/projects/{fundingId}", handlers.RemoveProjectFromGrant).Methods("DELETE")
	r.HandleFunc("/projects/{id}/grants", handlers.GetProjectGrants).Methods("GET")

	// --- RBAC ---
	r.HandleFunc("/roles", handlers.GetRoles).Methods("GET")
	r.HandleFunc("/permissions", handlers.GetPermissions).Methods("GET")
	r.HandleFunc("/projects/{id}/members/{userID}/role", handlers.AssignProjectRoleHandler).Methods("POST")
	r.HandleFunc("/projects/{id}/my-permissions", handlers.GetMyProjectPermissions).Methods("GET")
	r.HandleFunc("/projects/{id}/audit-log", handlers.GetProjectAuditLog).Methods("GET")

	// --- Отчеты ---
	r.HandleFunc("/projects/{id}/report", handlers.GenerateGostReport).Methods("GET")

	return r
}

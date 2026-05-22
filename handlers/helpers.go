package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"project-MVP/middleware"
	"project-MVP/services"
)

// sanitizeContentDisposition удаляет символы перевода строки из имени файла,
// чтобы предотвратить CRLF-инъекцию в HTTP-заголовке Content-Disposition.
func sanitizeContentDisposition(name string) string {
	name = strings.ReplaceAll(name, "\r", "")
	name = strings.ReplaceAll(name, "\n", "")
	name = strings.ReplaceAll(name, "\x00", "")
	return name
}

// GetUserID извлекает ID текущего пользователя из контекста запроса
func GetUserID(r *http.Request) (int, bool) {
	return middleware.GetUserID(r)
}

// RequireAuth требует наличия авторизации. При отсутствии пишет 401 и возвращает false.
func RequireAuth(w http.ResponseWriter, r *http.Request) (int, bool) {
	userID, ok := GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return 0, false
	}
	return userID, true
}

// RequirePermission проверяет право доступа к проекту. При отказе пишет 403 и возвращает false.
func RequirePermission(w http.ResponseWriter, r *http.Request, projectID int, permissionCode string) bool {
	userID, ok := RequireAuth(w, r)
	if !ok {
		return false
	}
	if err := services.CheckPermission(userID, projectID, permissionCode); err != nil {
		http.Error(w, "insufficient permissions", http.StatusForbidden)
		return false
	}
	return true
}

// RequireAdmin проверяет, что пользователь — системный администратор.
func RequireAdmin(w http.ResponseWriter, r *http.Request) (int, bool) {
	userID, ok := RequireAuth(w, r)
	if !ok {
		return 0, false
	}
	if !services.IsAdmin(userID) {
		http.Error(w, "admin access required", http.StatusForbidden)
		return 0, false
	}
	return userID, true
}

// GetProjectIDFromVars извлекает project_id из URL vars
func GetProjectIDFromVars(r *http.Request, key string) (int, bool) {
	vars := r.URL.Query()
	// Для gorilla/mux vars нужен mux.Vars, но здесь оставляем простой парсинг
	// Хендлеры сами вызывают mux.Vars(r) и передают projectID
	_ = vars
	return 0, false
}

// AtoiParam парсит строковый параметр в int
func AtoiParam(s string) (int, bool) {
	v, err := strconv.Atoi(s)
	if err != nil || v <= 0 {
		return 0, false
	}
	return v, true
}

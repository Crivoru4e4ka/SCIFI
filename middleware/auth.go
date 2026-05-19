package middleware

import (
	"context"
	"net/http"
	"strconv"
)

type ContextKey string

const userIDKey ContextKey = "userID"

// AuthMiddleware извлекает user_id из session-cookie и помещает его в контекст запроса.
// Не прерывает цепочку — позволяет публичным endpoint'ам работать без авторизации.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		userID, err := strconv.Atoi(cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID возвращает user_id из контекста запроса
func GetUserID(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(userIDKey).(int)
	return userID, ok
}

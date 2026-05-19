package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAuthMiddleware_ValidCookie проверяет, что middleware извлекает user_id
// из валидной session-cookie и помещает его в контекст запроса.
func TestAuthMiddleware_ValidCookie(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserID(r)
		assert.True(t, ok, "userID должен быть в контексте")
		assert.Equal(t, 42, userID, "userID должен совпадать со значением из cookie")
	})

	middleware := AuthMiddleware(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "42"})

	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)
}

// TestAuthMiddleware_NoCookie проверяет, что запрос без cookie
// проходит дальше, но userID отсутствует в контексте.
// Это важно для публичных endpoint'ов.
func TestAuthMiddleware_NoCookie(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := GetUserID(r)
		assert.False(t, ok, "userID не должен быть в контексте при отсутствии cookie")
	})

	middleware := AuthMiddleware(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)
}

// TestAuthMiddleware_InvalidCookieValue проверяет обработку cookie,
// значение которой не является валидным целым числом.
// Middleware не должен падать и должен пропустить запрос дальше без userID.
func TestAuthMiddleware_InvalidCookieValue(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := GetUserID(r)
		assert.False(t, ok, "userID не должен быть в контексте при невалидном значении cookie")
	})

	middleware := AuthMiddleware(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "not-a-number"})

	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)
}

// TestGetUserID_ExistingKey проверяет извлечение userID из контекста,
// когда ключ был установлен вручную.
func TestGetUserID_ExistingKey(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, userIDKey, 100)
	req = req.WithContext(ctx)

	userID, ok := GetUserID(req)
	assert.True(t, ok)
	assert.Equal(t, 100, userID)
}

// TestGetUserID_MissingKey проверяет, что GetUserID корректно
// сообщает об отсутствии ключа в контексте.
func TestGetUserID_MissingKey(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	userID, ok := GetUserID(req)
	assert.False(t, ok)
	assert.Equal(t, 0, userID)
}

// TestGetUserID_WrongType проверяет защиту от некорректного типа
// значения, сохраненного в контексте под ключом userIDKey.
func TestGetUserID_WrongType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, userIDKey, "not-an-int")
	req = req.WithContext(ctx)

	userID, ok := GetUserID(req)
	assert.False(t, ok)
	assert.Equal(t, 0, userID)
}

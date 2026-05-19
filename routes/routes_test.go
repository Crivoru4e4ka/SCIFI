package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNoCacheMiddleware_SetHeaders проверяет, что middleware корректно
// устанавливает заголовки, запрещающие кэширование.
// Это критически важно для безопасности личного кабинета.
func TestNoCacheMiddleware_SetHeaders(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := NoCacheMiddleware(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "no-cache, no-store, must-revalidate", rr.Header().Get("Cache-Control"))
	assert.Equal(t, "no-cache", rr.Header().Get("Pragma"))
	assert.Equal(t, "0", rr.Header().Get("Expires"))
}

// TestNoCacheMiddleware_NextCalled проверяет, что middleware
// передает управление следующему обработчику в цепочке.
func TestNoCacheMiddleware_NextCalled(t *testing.T) {
	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	middleware := NoCacheMiddleware(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)

	assert.True(t, called, "следующий обработчик должен быть вызван")
}

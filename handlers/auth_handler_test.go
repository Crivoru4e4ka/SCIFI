package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLogout_ClearsCookie проверяет, что хендлер Logout корректно
// удаляет сессионную cookie и перенаправляет на страницу логина.
func TestLogout_ClearsCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	rr := httptest.NewRecorder()

	Logout(rr, req)

	// Проверяем редирект на /login
	assert.Equal(t, http.StatusSeeOther, rr.Code)
	assert.Equal(t, "/login", rr.Header().Get("Location"))

	// Проверяем, что cookie удалена (MaxAge = -1, пустое значение)
	cookies := rr.Result().Cookies()
	require.Len(t, cookies, 1)
	cookie := cookies[0]

	assert.Equal(t, "session", cookie.Name)
	assert.Equal(t, "", cookie.Value)
	assert.Equal(t, -1, cookie.MaxAge)
	assert.True(t, cookie.Expires.IsZero() || cookie.Expires.Year() < 1971, "cookie должна иметь нулевое/истекшее время")
}

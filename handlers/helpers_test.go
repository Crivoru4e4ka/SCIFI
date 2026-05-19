package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"project-MVP/middleware"

	"github.com/stretchr/testify/assert"
)

// TestAtoiParam_Valid проверяет парсинг корректных положительных чисел.
func TestAtoiParam_Valid(t *testing.T) {
	val, ok := AtoiParam("42")
	assert.True(t, ok)
	assert.Equal(t, 42, val)
}

// TestAtoiParam_Zero проверяет, что ноль считается невалидным.
// Это важно для ID сущностей, которые начинаются с 1.
func TestAtoiParam_Zero(t *testing.T) {
	val, ok := AtoiParam("0")
	assert.False(t, ok)
	assert.Equal(t, 0, val)
}

// TestAtoiParam_Negative проверяет, что отрицательные числа отклоняются.
func TestAtoiParam_Negative(t *testing.T) {
	val, ok := AtoiParam("-5")
	assert.False(t, ok)
	assert.Equal(t, 0, val)
}

// TestAtoiParam_InvalidString проверяет обработку нечисловых строк.
func TestAtoiParam_InvalidString(t *testing.T) {
	val, ok := AtoiParam("abc")
	assert.False(t, ok)
	assert.Equal(t, 0, val)
}

// TestAtoiParam_Empty проверяет обработку пустой строки.
func TestAtoiParam_Empty(t *testing.T) {
	val, ok := AtoiParam("")
	assert.False(t, ok)
	assert.Equal(t, 0, val)
}

// TestRequireAuth_Authorized проверяет успешное извлечение userID
// из контекста запроса.
func TestRequireAuth_Authorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(req.Context(), middleware.ContextKey("userID"), 7)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	userID, ok := RequireAuth(rr, req)

	assert.True(t, ok)
	assert.Equal(t, 7, userID)
	assert.Equal(t, http.StatusOK, rr.Code)
}

// TestRequireAuth_Unauthorized проверяет, что при отсутствии userID
// в контексте возвращается HTTP 401 Unauthorized.
func TestRequireAuth_Unauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	userID, ok := RequireAuth(rr, req)

	assert.False(t, ok)
	assert.Equal(t, 0, userID)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

// TestGetUserID_Delegation проверяет, что хендлерский GetUserID
// делегирует вызов middleware.GetUserID.
func TestGetUserID_Delegation(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(req.Context(), middleware.ContextKey("userID"), 99)
	req = req.WithContext(ctx)

	userID, ok := GetUserID(req)
	assert.True(t, ok)
	assert.Equal(t, 99, userID)
}

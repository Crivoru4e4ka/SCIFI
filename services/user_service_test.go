package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRoleIsValid проверяет валидацию системных ролей пользователей.
// Роли чувствительны к регистру и пробелам — функция нормализует вход.
func TestRoleIsValid(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected bool
	}{
		{"admin lowercase", "admin", true},
		{"user lowercase", "user", true},
		{"guest lowercase", "guest", true},
		{"admin uppercase", "ADMIN", true},
		{"mixed case", "Admin", true},
		{"with spaces", "  admin  ", true},
		{"empty string", "", false},
		{"unknown role", "superuser", false},
		{"whitespace only", "   ", false},
		{"similar but invalid", "administrator", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := roleIsValid(tt.role)
			assert.Equal(t, tt.expected, result, "role: %q", tt.role)
		})
	}
}

// TestAllowedRolesMap проверяет, что карта допустимых ролей
// содержит только ожидаемые значения. Это защита от случайной
// модификации глобальной переменной.
func TestAllowedRolesMap(t *testing.T) {
	assert.Len(t, allowedRoles, 3)
	assert.Contains(t, allowedRoles, "admin")
	assert.Contains(t, allowedRoles, "user")
	assert.Contains(t, allowedRoles, "guest")
	assert.NotContains(t, allowedRoles, "superuser")
}

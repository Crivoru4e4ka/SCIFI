package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsValidProjectRole_CacheHit проверяет валидацию проектных ролей
// при наличии данных в кэше. Мы напрямую заполняем кэш, минуя БД.
func TestIsValidProjectRole_CacheHit(t *testing.T) {
	// Подготавливаем кэш вручную (white-box)
	projectRoleNamesMutex.Lock()
	oldCache := projectRoleNamesCache
	projectRoleNamesCache = map[string]bool{
		"project_lead": true,
		"researcher":   true,
		"viewer":       true,
	}
	projectRoleNamesMutex.Unlock()

	// Восстанавливаем после теста
	defer func() {
		projectRoleNamesMutex.Lock()
		projectRoleNamesCache = oldCache
		projectRoleNamesMutex.Unlock()
	}()

	tests := []struct {
		name     string
		role     string
		expected bool
	}{
		{"valid project_lead", "project_lead", true},
		{"valid researcher", "researcher", true},
		{"valid uppercase", "PROJECT_LEAD", true},
		{"valid with spaces", "  viewer  ", true},
		{"invalid role", "supervisor", false},
		{"empty string", "", false},
		{"whitespace only", "   ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidProjectRole(tt.role)
			assert.Equal(t, tt.expected, result, "role: %q", tt.role)
		})
	}
}

// TestGetProjectRoleNames проверяет извлечение списка проектных ролей из кэша.
func TestGetProjectRoleNames(t *testing.T) {
	projectRoleNamesMutex.Lock()
	oldCache := projectRoleNamesCache
	projectRoleNamesCache = map[string]bool{
		"analyst": true,
		"viewer":  true,
	}
	projectRoleNamesMutex.Unlock()

	defer func() {
		projectRoleNamesMutex.Lock()
		projectRoleNamesCache = oldCache
		projectRoleNamesMutex.Unlock()
	}()

	names := GetProjectRoleNames()
	assert.Len(t, names, 2)
	assert.Contains(t, names, "analyst")
	assert.Contains(t, names, "viewer")
}

// TestRolePermissionMappingsIntegrity проверяет, что матрица доступа
// содержит только роли, определенные в projectRoles и systemRoles.
// Это важно для консистентности RBAC-системы.
func TestRolePermissionMappingsIntegrity(t *testing.T) {
	validRoles := map[string]bool{
		"admin": true, "user": true, "guest": true,
	}
	for _, r := range projectRoles {
		validRoles[r.Name] = true
	}

	for roleName := range rolePermissionMappings {
		assert.Contains(t, validRoles, roleName, "роль %s должна быть определена в системе", roleName)
	}
}

// TestPermissionDefinitionsUniqueness проверяет, что все коды прав уникальны.
// Дублирующиеся коды приведут к конфликтам при вставке в БД.
func TestPermissionDefinitionsUniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for _, p := range permissionDefinitions {
		assert.False(t, seen[p.Code], "код права %s дублируется", p.Code)
		seen[p.Code] = true
	}
}

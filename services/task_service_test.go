package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTaskStatusIsValid проверяет валидацию статусов задач.
// Допустимые статусы: todo, in_progress, review, done.
func TestTaskStatusIsValid(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{"todo", "todo", true},
		{"in_progress", "in_progress", true},
		{"review", "review", true},
		{"done", "done", true},
		{"TODO uppercase", "TODO", true},
		{"In_Progress mixed", "In_Progress", true},
		{"with spaces", "  done  ", true},
		{"empty string", "", false},
		{"unknown status", "blocked", false},
		{"whitespace only", "   ", false},
		{"archived", "archived", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := taskStatusIsValid(tt.status)
			assert.Equal(t, tt.expected, result, "status: %q", tt.status)
		})
	}
}

// TestAllowedTaskStatusesMap проверяет целостность глобальной карты
// допустимых статусов задач.
func TestAllowedTaskStatusesMap(t *testing.T) {
	assert.Len(t, allowedTaskStatuses, 4)
	assert.Contains(t, allowedTaskStatuses, "todo")
	assert.Contains(t, allowedTaskStatuses, "in_progress")
	assert.Contains(t, allowedTaskStatuses, "review")
	assert.Contains(t, allowedTaskStatuses, "done")
}

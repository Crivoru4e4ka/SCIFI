package services

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGrantTypeIsValid проверяет валидацию типов грантов.
func TestGrantTypeIsValid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"state lowercase", "state", true},
		{"university", "university", true},
		{"international", "international", true},
		{"corporate", "corporate", true},
		{"internal", "internal", true},
		{"uppercase", "STATE", true},
		{"mixed case", "University", true},
		{"with spaces", "  corporate  ", true},
		{"empty", "", false},
		{"unknown", "federal", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, grantTypeIsValid(tt.input))
		})
	}
}

// TestGrantStatusIsValid проверяет валидацию статусов грантов.
func TestGrantStatusIsValid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"draft", "draft", true},
		{"submitted", "submitted", true},
		{"under_review", "under_review", true},
		{"approved", "approved", true},
		{"rejected", "rejected", true},
		{"active", "active", true},
		{"completed", "completed", true},
		{"suspended", "suspended", true},
		{" uppercase", "ACTIVE", true},
		{"empty", "", false},
		{"unknown", "cancelled", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, grantStatusIsValid(tt.input))
		})
	}
}

// TestParseDatePtr_Valid проверяет парсинг корректной даты.
func TestParseDatePtr_Valid(t *testing.T) {
	s := "2024-05-20"
	result, err := parseDatePtr(&s)

	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2024, result.Year())
	assert.Equal(t, time.May, result.Month())
	assert.Equal(t, 20, result.Day())
}

// TestParseDatePtr_Nil проверяет обработку nil-указателя.
func TestParseDatePtr_Nil(t *testing.T) {
	result, err := parseDatePtr(nil)

	assert.NoError(t, err)
	assert.Nil(t, result)
}

// TestParseDatePtr_Empty проверяет обработку пустой строки.
func TestParseDatePtr_Empty(t *testing.T) {
	s := ""
	result, err := parseDatePtr(&s)

	assert.NoError(t, err)
	assert.Nil(t, result)
}

// TestParseDatePtr_Invalid проверяет обработку невалидной даты.
func TestParseDatePtr_Invalid(t *testing.T) {
	s := "not-a-date"
	result, err := parseDatePtr(&s)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestNullTimeStr_Valid проверяет конвертацию валидного sql.NullTime.
func TestNullTimeStr_Valid(t *testing.T) {
	nt := sql.NullTime{Valid: true, Time: time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC)}
	result := nullTimeStr(nt)

	assert.NotNil(t, result)
	assert.Equal(t, "2024-05-20", *result)
}

// TestNullTimeStr_Invalid проверяет конвертацию невалидного sql.NullTime.
func TestNullTimeStr_Invalid(t *testing.T) {
	nt := sql.NullTime{Valid: false}
	result := nullTimeStr(nt)

	assert.Nil(t, result)
}

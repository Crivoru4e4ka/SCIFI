package services

import (
	"testing"
	"time"

	"project-MVP/models"

	"github.com/stretchr/testify/assert"
)

// TestValidateAndNormalizeProject_Success проверяет успешную
// нормализацию проекта с установкой значений по умолчанию.
func TestValidateAndNormalizeProject_Success(t *testing.T) {
	teamID := 5
	p := models.Project{
		Name:          "  AI Research  ",
		Description:   "  Description  ",
		CreatedBy:     1,
		ExecutionType: "  TEAM  ",
		TeamId:        &teamID,
	}

	err := ValidateAndNormalizeProject(&p)

	assert.NoError(t, err)
	assert.Equal(t, "AI Research", p.Name)
	assert.Equal(t, "Description", p.Description)
	assert.Equal(t, "team", p.ExecutionType)
	assert.Equal(t, "active", p.Status)
	assert.Equal(t, "closed", p.Visibility)
	assert.False(t, p.StartDate.IsZero())
}

// TestValidateAndNormalizeProject_EmptyName проверяет отказ
// при пустом названии проекта.
func TestValidateAndNormalizeProject_EmptyName(t *testing.T) {
	p := models.Project{
		Name:      "   ",
		CreatedBy: 1,
	}

	err := ValidateAndNormalizeProject(&p)

	assert.EqualError(t, err, "project name is required")
}

// TestValidateAndNormalizeProject_InvalidCreatedBy проверяет отказ
// при отсутствии создателя проекта.
func TestValidateAndNormalizeProject_InvalidCreatedBy(t *testing.T) {
	p := models.Project{
		Name:      "Project",
		CreatedBy: 0,
	}

	err := ValidateAndNormalizeProject(&p)

	assert.EqualError(t, err, "created_by is required")
}

// TestValidateAndNormalizeProject_InvalidExecutionType проверяет
// отказ при недопустимом типе исполнения.
func TestValidateAndNormalizeProject_InvalidExecutionType(t *testing.T) {
	p := models.Project{
		Name:          "Project",
		CreatedBy:     1,
		ExecutionType: "invalid",
	}

	err := ValidateAndNormalizeProject(&p)

	assert.EqualError(t, err, "invalid execution type")
}

// TestValidateAndNormalizeProject_TeamWithoutTeamID проверяет
// обязательность team_id для проектов типа "team".
func TestValidateAndNormalizeProject_TeamWithoutTeamID(t *testing.T) {
	p := models.Project{
		Name:          "Project",
		CreatedBy:     1,
		ExecutionType: "team",
		TeamId:        nil,
	}

	err := ValidateAndNormalizeProject(&p)

	assert.EqualError(t, err, "team_id is required for team execution type")
}

// TestValidateAndNormalizeProject_TeamWithZeroTeamID проверяет,
// что team_id = 0 считается невалидным.
func TestValidateAndNormalizeProject_TeamWithZeroTeamID(t *testing.T) {
	teamID := 0
	p := models.Project{
		Name:          "Project",
		CreatedBy:     1,
		ExecutionType: "team",
		TeamId:        &teamID,
	}

	err := ValidateAndNormalizeProject(&p)

	assert.EqualError(t, err, "team_id is required for team execution type")
}

// TestValidateAndNormalizeProject_DefaultValues проверяет
// автоматическую установку значений по умолчанию.
func TestValidateAndNormalizeProject_DefaultValues(t *testing.T) {
	p := models.Project{
		Name:      "Project",
		CreatedBy: 1,
	}

	err := ValidateAndNormalizeProject(&p)

	assert.NoError(t, err)
	assert.Equal(t, "active", p.Status)
	assert.Equal(t, "closed", p.Visibility)
	assert.Equal(t, "manual", p.ExecutionType)
	assert.WithinDuration(t, time.Now(), p.StartDate, time.Second)
}

// TestValidateAndNormalizeProject_PreservesExplicitValues проверяет,
// что явно заданные значения не перезаписываются дефолтными.
func TestValidateAndNormalizeProject_PreservesExplicitValues(t *testing.T) {
	future := time.Now().Add(24 * time.Hour)
	teamID := 10
	p := models.Project{
		Name:          "Project",
		CreatedBy:     1,
		Status:        "archived",
		Visibility:    "open",
		ExecutionType: "team",
		TeamId:        &teamID,
		StartDate:     future,
	}

	err := ValidateAndNormalizeProject(&p)

	assert.NoError(t, err)
	assert.Equal(t, "archived", p.Status)
	assert.Equal(t, "open", p.Visibility)
	assert.Equal(t, "team", p.ExecutionType)
	assert.Equal(t, future, p.StartDate)
	assert.Equal(t, 10, *p.TeamId)
}

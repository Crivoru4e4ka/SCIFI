package services

import (
	"os"
	"testing"

	"project-MVP/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateProjectExcel_EmptyTasks проверяет генерацию Excel
// для проекта без задач. Должен создаваться валидный файл с заголовками.
func TestGenerateProjectExcel_EmptyTasks(t *testing.T) {
	project := models.Project{
		Id:   1,
		Name: "Test Project",
		Key:  "TP",
	}

	buf, err := GenerateProjectExcel(project, []models.Task{})
	require.NoError(t, err)
	assert.NotNil(t, buf)
	assert.Greater(t, buf.Len(), 0)
}

// TestGenerateProjectExcel_WithTasks проверяет генерацию Excel
// с задачами различных типов, включая nil-указатели.
func TestGenerateProjectExcel_WithTasks(t *testing.T) {
	project := models.Project{
		Id:   1,
		Name: "Research Project",
		Key:  "RP",
	}

	taskType := "experiment"
	tasks := []models.Task{
		{
			Id:       1,
			TaskNum:  1,
			Title:    "Task One",
			Status:   "done",
			Priority: "high",
			Type:     &taskType,
			ResearchContribution: "confirmed",
			ResearchMethod:       "simulation",
			DOI:                  "10.1234/test",
			Conclusion:           "Positive result",
		},
		{
			Id:       2,
			TaskNum:  2,
			Title:    "Task Two",
			Status:   "todo",
			Priority: "medium",
			Type:     nil, // nil Type — важный edge case
		},
	}

	buf, err := GenerateProjectExcel(project, tasks)
	require.NoError(t, err)
	assert.NotNil(t, buf)
	assert.Greater(t, buf.Len(), 0)
}

// TestGenerateProjectPDF_EmptyTasks проверяет генерацию PDF
// для проекта без задач. Пропускается, если шрифты отсутствуют в ФС.
func TestGenerateProjectPDF_EmptyTasks(t *testing.T) {
	if _, err := os.Stat("static/fonts/Arial.ttf"); os.IsNotExist(err) {
		t.Skip("шрифты PDF недоступны, пропускаем тест")
	}

	project := models.Project{
		Id:             1,
		Name:           "Test Project",
		Key:            "TP",
		MainHypothesis: "H1",
		ResearchGoal:   "Goal",
	}

	data, err := GenerateProjectPDF(project, []models.Task{})
	require.NoError(t, err)
	assert.NotNil(t, data)
	assert.Greater(t, len(data), 0)
	// Простая проверка PDF-заголовка (magic bytes %PDF)
	assert.Equal(t, "%PDF", string(data[:4]))
}

// TestGenerateProjectPDF_WithTasks проверяет генерацию PDF
// с задачами. Пропускается, если шрифты отсутствуют в ФС.
func TestGenerateProjectPDF_WithTasks(t *testing.T) {
	if _, err := os.Stat("static/fonts/Arial.ttf"); os.IsNotExist(err) {
		t.Skip("шрифты PDF недоступны, пропускаем тест")
	}

	project := models.Project{
		Id:             1,
		Name:           "Research Project",
		Key:            "RP",
		MainHypothesis: "Main hypothesis",
		ResearchGoal:   "Research goal",
	}

	tasks := []models.Task{
		{
			Id:      1,
			TaskNum: 1,
			Title:   "Experiment 1",
			Status:  "done",
		},
		{
			Id:      2,
			TaskNum: 2,
			Title:   "Analysis",
			Status:  "in_progress",
		},
	}

	data, err := GenerateProjectPDF(project, tasks)
	require.NoError(t, err)
	assert.NotNil(t, data)
	assert.Greater(t, len(data), 0)
	assert.Equal(t, "%PDF", string(data[:4]))
}

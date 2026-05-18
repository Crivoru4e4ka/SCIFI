package services

import (
	"bytes"
	"fmt"
	"project-MVP/models"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

// GenerateProjectExcel создает Excel файл со всеми задачами проекта
func GenerateProjectExcel(project models.Project, tasks []models.Task) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheet := "Задачи"
	f.SetSheetName("Sheet1", sheet)

	// Заголовки
	headers := []string{"Ключ", "Название", "Тип", "Статус", "Приоритет", "Вклад в гипотезу", "Метод", "DOI", "Вывод"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Данные
	for i, t := range tasks {
		row := i + 2
		taskKey := fmt.Sprintf("%s-%d", project.Key, t.TaskNum)

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), taskKey)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), t.Title)
		// Проверяем, что указатель не пустой, и берем само значение
		taskType := ""
		if t.Type != nil {
			taskType = *t.Type
		}
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), taskType)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), t.Status)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), t.Priority)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), t.ResearchContribution)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), t.ResearchMethod)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), t.DOI)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), t.Conclusion)
	}

	return f.WriteToBuffer()
}

// GenerateProjectPDF создает PDF отчет
func GenerateProjectPDF(project models.Project, tasks []models.Task) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("Arial", "", "static/fonts/Arial.ttf")
	pdf.AddUTF8Font("Arial", "B", "static/fonts/Arial-Bold.ttf")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, project.Name)
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(190, 5, "Основная гипотеза: "+project.MainHypothesis, "", "L", false)
	pdf.Ln(5)
	pdf.MultiCell(190, 5, "Цель исследования: "+project.ResearchGoal, "", "L", false)
	pdf.Ln(10)

	// Заголовки таблицы
	pdf.SetFont("Arial", "B", 12)
	// Параметры CellFormat: ширина, высота, текст, рамка, переход на новую строку, выравнивание, заливка, ссылка, текст ссылки
	pdf.CellFormat(40, 7, "Key", "1", 0, "C", false, 0, "")
	pdf.CellFormat(100, 7, "Title", "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 7, "Status", "1", 0, "C", false, 0, "")
	pdf.Ln(7)

	// Данные таблицы
	pdf.SetFont("Arial", "", 10)
	for _, t := range tasks {
		taskKey := fmt.Sprintf("%s-%d", project.Key, t.TaskNum)

		pdf.CellFormat(40, 6, taskKey, "1", 0, "L", false, 0, "")
		pdf.CellFormat(100, 6, t.Title, "1", 0, "L", false, 0, "")
		pdf.CellFormat(50, 6, t.Status, "1", 0, "L", false, 0, "")
		pdf.Ln(6)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}

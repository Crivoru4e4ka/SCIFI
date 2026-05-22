package handlers

import (
	"fmt"
	"net/http"
	"project-MVP/services"
	"strconv"

	"github.com/gorilla/mux"
)

// ExportProjectExcelHandler godoc
// @Summary Выгрузить задачи в Excel
// @Description Генерирует XLSX файл, содержащий подробный список всех задач проекта со всеми научными параметрами
// @Tags export
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param id path int true "Project ID"
// @Success 200 {file} binary "Файл отчета .xlsx"
// @Failure 403 {string} string "insufficient permissions"
// @Failure 500 {string} string "internal error"
// @Router /projects/{id}/export/excel [get]
func ExportProjectExcelHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, _ := strconv.Atoi(vars["id"])

	// 1. Проверка прав
	if !RequirePermission(w, r, projectID, "project.view") {
		return
	}

	// 2. Получение данных
	project, err := services.GetProjectByID(projectID)
	if err != nil {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}
	tasks, err := services.GetTasksByProject(projectID)
	if err != nil {
		http.Error(w, "failed to load tasks", http.StatusInternalServerError)
		return
	}

	// 3. Генерация
	buffer, err := services.GenerateProjectExcel(project, tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Отправка файла (с защитой от CRLF-инъекции в заголовке)
	fileName := sanitizeContentDisposition(fmt.Sprintf("Report_%s.xlsx", project.Key))
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Write(buffer.Bytes())
}

// ExportProjectPDFHandler godoc
// @Summary Выгрузить отчет в PDF
// @Description Формирует PDF документ с краткой информацией о проекте, его целях, гипотезах и таблицей задач
// @Tags export
// @Produce application/pdf
// @Param id path int true "Project ID"
// @Success 200 {file} binary "Файл отчета .pdf"
// @Failure 403 {string} string "insufficient permissions"
// @Failure 500 {string} string "internal error"
// @Router /projects/{id}/export/pdf [get]
func ExportProjectPDFHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, _ := strconv.Atoi(vars["id"])

	if !RequirePermission(w, r, projectID, "project.view") {
		return
	}

	project, err := services.GetProjectByID(projectID)
	if err != nil {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}
	tasks, err := services.GetTasksByProject(projectID)
	if err != nil {
		http.Error(w, "failed to load tasks", http.StatusInternalServerError)
		return
	}

	pdfBytes, err := services.GenerateProjectPDF(project, tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileName := sanitizeContentDisposition(fmt.Sprintf("Report_%s.pdf", project.Key))
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Write(pdfBytes)
}

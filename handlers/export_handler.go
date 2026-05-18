package handlers

import (
	"fmt"
	"net/http"
	"project-MVP/services"
	"strconv"

	"github.com/gorilla/mux"
)

func ExportProjectExcelHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, _ := strconv.Atoi(vars["id"])

	// 1. Проверка прав
	if !RequirePermission(w, r, projectID, "project.view") {
		return
	}

	// 2. Получение данных
	project, _ := services.GetProjectByID(projectID)
	tasks, _ := services.GetTasksByProject(projectID)

	// 3. Генерация
	buffer, err := services.GenerateProjectExcel(project, tasks)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 4. Отправка файла
	fileName := fmt.Sprintf("Report_%s.xlsx", project.Key)
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Write(buffer.Bytes())
}

func ExportProjectPDFHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, _ := strconv.Atoi(vars["id"])

	if !RequirePermission(w, r, projectID, "project.view") {
		return
	}

	project, _ := services.GetProjectByID(projectID)
	tasks, _ := services.GetTasksByProject(projectID)

	pdfBytes, err := services.GenerateProjectPDF(project, tasks)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=Report_%s.pdf", project.Key))
	w.Write(pdfBytes)
}

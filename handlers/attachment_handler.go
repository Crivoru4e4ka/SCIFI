package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"project-MVP/models"
	"project-MVP/services"
	"strconv"

	"github.com/gorilla/mux"
)

func GetProjectAttachmentsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, _ := strconv.Atoi(vars["id"])

	if !RequirePermission(w, r, projectID, "project.view") {
		return
	}

	files, err := services.GetProjectAttachments(projectID)
	if err != nil {
		http.Error(w, "Ошибка при получении файлов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, _ := strconv.Atoi(vars["id"])

	// Получаем ID текущего пользователя из контекста
	userID, ok := RequireAuth(w, r)
	if !ok {
		return
	}

	// Получаем задачу для проверки project_id
	task, err := services.GetTaskByID(taskID)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	if !RequirePermission(w, r, task.ProjectId, "attachment.upload") {
		return
	}

	// Читаем файл
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Создаем папку, если её вдруг нет
	uploadDir := "./static/uploads"
	os.MkdirAll(uploadDir, os.ModePerm)

	// Формируем путь и сохраняем файл
	filePath := uploadDir + "/" + header.Filename
	out, err := os.Create(filePath)
	if err != nil {
		log.Printf("Ошибка при создании файла: %v", err)
		http.Error(w, "Не удалось сохранить файл на диске", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	io.Copy(out, file)

	// Записываем в базу
	attachment := models.Attachment{
		TaskId:   taskID,
		UserId:   userID,
		FileName: header.Filename,
		FileUrl:  "/static/uploads/" + header.Filename,
	}

	if err := services.CreateAttachment(&attachment); err != nil {
		log.Printf("Ошибка БД: %v", err)
		http.Error(w, "Ошибка записи в базу данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attachment)
}

func GenerateGostReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	projectID, _ := strconv.Atoi(idStr)

	if !RequirePermission(w, r, projectID, "project.view") {
		return
	}

	reportText, err := services.GetProjectReportData(projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=report.txt")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(reportText))
}

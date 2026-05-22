package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"project-MVP/models"
	"project-MVP/services"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// generateSafeFilename создаёт безопасное уникальное имя файла и санитизирует оригинальное имя.
func generateSafeFilename(original string) (sanitizedOriginal string, physicalName string) {
	// Санитизируем оригинальное имя: убираем path-разделители и traversal
	safe := filepath.Base(original)
	safe = strings.ReplaceAll(safe, "..", "")
	safe = strings.ReplaceAll(safe, "/", "")
	safe = strings.ReplaceAll(safe, "\\", "")
	if safe == "" || safe == "." {
		safe = "upload"
	}

	// Генерируем уникальное физическое имя
	ext := filepath.Ext(safe)
	b := make([]byte, 8)
	rand.Read(b)
	physical := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), hex.EncodeToString(b), ext)
	return safe, physical
}

// GetProjectAttachmentsHandler godoc
// @Summary Получить все вложения проекта
// @Description Возвращает список всех файлов, прикрепленных к задачам внутри указанного проекта
// @Tags attachments
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {array} models.Attachment "Список файлов"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "insufficient permissions"
// @Failure 500 {string} string "Ошибка при получении файлов"
// @Router /projects/{id}/attachments [get]
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

// UploadFileHandler godoc
// @Summary Загрузить файл к задаче
// @Description Сохраняет файл на сервере и привязывает его к конкретной научной задаче
// @Tags attachments
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Task ID"
// @Param file formData file true "Выбрать файл для загрузки"
// @Success 200 {object} models.Attachment "Данные о сохраненном файле"
// @Failure 400 {string} string "Ошибка чтения файла"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 404 {string} string "task not found"
// @Failure 500 {string} string "Ошибка записи"
// @Router /tasks/{id}/attachments [post]
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

	// Санитизируем имя файла и генерируем безопасное физическое имя (защита от Path Traversal)
	safeName, physicalName := generateSafeFilename(header.Filename)
	filePath := filepath.Join(uploadDir, physicalName)

	// Гарантируем, что финальный путь остаётся внутри uploadDir
	absUploadDir, _ := filepath.Abs(uploadDir)
	absFilePath, _ := filepath.Abs(filePath)
	if !strings.HasPrefix(absFilePath, absUploadDir+string(filepath.Separator)) {
		http.Error(w, "Недопустимое имя файла", http.StatusBadRequest)
		return
	}

	out, err := os.Create(filePath)
	if err != nil {
		log.Printf("Ошибка при создании файла: %v", err)
		http.Error(w, "Не удалось сохранить файл на диске", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	io.Copy(out, file)

	// Записываем в базу (сохраняем оригинальное санитизированное имя для отображения)
	attachment := models.Attachment{
		TaskId:   taskID,
		UserId:   userID,
		FileName: safeName,
		FileUrl:  "/static/uploads/" + physicalName,
	}

	if err := services.CreateAttachment(&attachment); err != nil {
		log.Printf("Ошибка БД: %v", err)
		http.Error(w, "Ошибка записи в базу данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attachment)
}

// GenerateGostReport godoc
// @Summary Генерировать отчет по ГОСТ
// @Description Генерирует отчет по ГОСТ для указанного проекта
// @Tags reports
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {string} string "Отчет в формате TXT"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "insufficient permissions"
// @Failure 500 {string} string "Ошибка при генерации отчета"
// @Router /projects/{id}/gost-report [get]
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

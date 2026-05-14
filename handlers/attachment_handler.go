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

	// 1. Получаем ID текущего пользователя из сессии
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "Нужна авторизация", http.StatusUnauthorized)
		return
	}
	userID, _ := strconv.Atoi(cookie.Value)

	// 2. Читаем файл
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 3. Создаем папку, если её вдруг нет (защита от ошибки 500)
	uploadDir := "./static/uploads"
	os.MkdirAll(uploadDir, os.ModePerm)

	// 4. Формируем путь и сохраняем файл
	filePath := uploadDir + "/" + header.Filename
	out, err := os.Create(filePath)
	if err != nil {
		log.Printf("Ошибка при создании файла: %v", err) // Это появится в консоли Go
		http.Error(w, "Не удалось сохранить файл на диске", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	io.Copy(out, file)

	// 5. Записываем в базу
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

	reportText, err := services.GetProjectReportData(projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=report.txt")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(reportText))
}

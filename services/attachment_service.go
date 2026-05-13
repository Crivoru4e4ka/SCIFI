package services

import (
	"log"
	"project-MVP/db"
	"project-MVP/models"
)

func GetProjectAttachments(projectId int) ([]models.Attachment, error) {
	// Выбираем вложения всех задач, которые принадлежат данному проекту
	query := `
        SELECT a.id, a.task_id, a.user_id, a.file_name, a.file_url, a.created_at 
        FROM attachments a
        JOIN tasks t ON a.task_id = t.id
        WHERE t.project_id = $1
    `
	rows, err := db.DB.Query(query, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []models.Attachment
	for rows.Next() {
		var a models.Attachment
		if err := rows.Scan(&a.Id, &a.TaskId, &a.UserId, &a.FileName, &a.FileUrl, &a.CreatedAt); err != nil {
			return nil, err
		}
		attachments = append(attachments, a)
	}
	return attachments, nil
}

func CreateAttachment(a *models.Attachment) error {
	// Это сообщение покажет, что функция ВООБЩЕ запустилась
	log.Printf("БД: Начинаем вставку в таблицу attachments. TaskID: %d, UserID: %d", a.TaskId, a.UserId)

	query := `
		INSERT INTO attachments (task_id, user_id, file_name, file_url, created_at) 
		VALUES ($1, $2, $3, $4, NOW()) 
		RETURNING id, created_at`

	err := db.DB.QueryRow(query, a.TaskId, a.UserId, a.FileName, a.FileUrl).Scan(&a.Id, &a.CreatedAt)

	if err != nil {
		// ВАЖНО: это сообщение ОБЯЗАТЕЛЬНО появится в терминале при ошибке
		log.Printf("!!! КРИТИЧЕСКАЯ ОШИБКА БД !!!: %v", err)
		return err
	}

	log.Printf("БД: Запись успешно создана! Новый ID вложения: %d", a.Id)
	return nil
}

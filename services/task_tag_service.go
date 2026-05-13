package services

import (
	"database/sql"
	"project-MVP/db"
	"project-MVP/models"
)

// AddTagToTask добавляет тег к задаче (создает тег, если его нет)
func AddTagToTask(taskID int, tag models.Tag) (models.Task, error) {
	// 1. Проверяем, существует ли задача (используем _, чтобы не было ошибки unused variable)
	_, err := GetTaskByID(taskID)
	if err != nil {
		return models.Task{}, err
	}

	// 2. Проверяем, существует ли тег с таким именем
	var existingTag models.Tag
	// Используем $1 вместо ?, если у вас PostgreSQL
	err = db.DB.QueryRow("SELECT id, name FROM tags WHERE name = $1", tag.Name).Scan(&existingTag.ID, &existingTag.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			// Если тега нет - создаем (вызывается из tag_service.go)
			// Убедитесь, что в модели Tag есть ProjectID, если он обязателен в БД
			existingTag, err = CreateTag(tag)
			if err != nil {
				return models.Task{}, err
			}
		} else {
			return models.Task{}, err
		}
	}

	// 3. Привязываем тег к задаче в таблице связей
	// INSERT INTO task_tags (task_id, tag_id) ...
	_, err = db.DB.Exec("INSERT INTO task_tags (task_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", taskID, existingTag.ID)
	if err != nil {
		return models.Task{}, err
	}

	// 4. Возвращаем обновленную задачу
	return GetTaskByID(taskID)
}

// GetTagsForTask — ЭТОЙ ФУНКЦИИ НЕ ХВАТАЛО
func GetTagsForTask(taskID int) ([]models.Tag, error) {
	query := `
		SELECT t.id, t.name, t.created_at 
		FROM tags t
		JOIN task_tags tt ON t.id = tt.tag_id
		WHERE tt.task_id = $1`

	rows, err := db.DB.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var t models.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, nil
}

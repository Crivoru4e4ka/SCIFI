package services

import (
	"project-MVP/db"
	"project-MVP/models"
)

func AddComment(c *models.Comment) (int, error) {
	var id int
	query := `
		INSERT INTO comments (entity_id, entity_type, user_id, parent_id, content)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err := db.DB.QueryRow(query, c.EntityId, c.EntityType, c.UserId, c.ParentId, c.Content).Scan(&id)
	return id, err
}

func GetCommentsByEntity(entityType string, entityId int) ([]models.Comment, error) {
	query := `
		SELECT c.id, c.entity_type, c.entity_id, c.user_id, u.full_name, 
		       c.parent_id, COALESCE(c.content, '') as content, 
               c.created_at, c.updated_at, c.deleted_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.entity_type = $1 AND c.entity_id = $2
		ORDER BY c.created_at ASC`

	rows, err := db.DB.Query(query, entityType, entityId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 1. Создаем список указателей. Это важно, чтобы работать с конкретными объектами в памяти.
	allComments := []*models.Comment{}
	for rows.Next() {
		c := &models.Comment{} // Создаем указатель
		err := rows.Scan(&c.Id, &c.EntityType, &c.EntityId, &c.UserId, &c.UserName,
			&c.ParentId, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt)
		if err != nil {
			return nil, err
		}

		if c.DeletedAt != nil {
			c.Content = "Комментарий удален"
		}

		// Инициализируем пустой слайс, чтобы фронтенд получил [] вместо null
		c.Replies = []models.Comment{}
		allComments = append(allComments, c)
	}

	// 2. Делаем карту (Map) для быстрого поиска родителей
	commentMap := make(map[int]*models.Comment)
	for _, c := range allComments {
		commentMap[c.Id] = c
	}

	// 3. Строим дерево
	finalRoots := []models.Comment{}
	for _, c := range allComments {
		if c.ParentId == nil {
			// Если нет родителя — это корень. Мы добавим его в результат позже.
			continue
		} else {
			// Если есть родитель — ищем его в карте и добавляем ответ К НЕМУ
			if parent, ok := commentMap[*c.ParentId]; ok {
				parent.Replies = append(parent.Replies, *c)
			}
		}
	}

	// 4. Собираем только корневые комментарии в финальный список
	// (теперь у них внутри уже лежат ответы, добавленные на шаге 3)
	for _, c := range allComments {
		if c.ParentId == nil {
			finalRoots = append(finalRoots, *c)
		}
	}

	return finalRoots, nil
}

func SoftDeleteComment(id int) error {
	query := `UPDATE comments SET deleted_at = NOW(), content = NULL WHERE id = $1`
	_, err := db.DB.Exec(query, id)
	return err
}

func GetCommentRaw(id int) (*models.Comment, error) {
	var c models.Comment
	err := db.DB.QueryRow("SELECT id, user_id, entity_id, entity_type FROM comments WHERE id = $1", id).
		Scan(&c.Id, &c.UserId, &c.EntityId, &c.EntityType)
	return &c, err
}

func UpdateComment(id int, content string) error {
	_, err := db.DB.Exec("UPDATE comments SET content = $1, updated_at = NOW() WHERE id = $2", content, id)
	return err
}

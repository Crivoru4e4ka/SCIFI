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
		       c.parent_id, COALESCE(c.content, '') as content, -- ДОБАВЛЕНО COALESCE
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

	all := []models.Comment{}
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(&c.Id, &c.EntityType, &c.EntityId, &c.UserId, &c.UserName,
			&c.ParentId, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt)
		if err != nil {
			return nil, err
		}
		if c.DeletedAt != nil {
			c.Content = "Комментарий удален"
		}
		all = append(all, c)
	}

	// Построение дерева (1 уровень вложенности)
	roots := []models.Comment{}
	commentMap := make(map[int]*models.Comment)

	for i := range all {
		if all[i].ParentId == nil {
			roots = append(roots, all[i])
			commentMap[all[i].Id] = &roots[len(roots)-1]
		}
	}

	for i := range all {
		if all[i].ParentId != nil {
			if parent, ok := commentMap[*all[i].ParentId]; ok {
				parent.Replies = append(parent.Replies, all[i])
			}
		}
	}
	return roots, nil
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

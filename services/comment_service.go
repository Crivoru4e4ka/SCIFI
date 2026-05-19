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

	return BuildCommentTree(allComments), nil
}

// BuildCommentTree строит дерево комментариев из плоского списка.
// Корневые комментарии (ParentId == nil) возвращаются на верхнем уровне,
// а ответы (replies) рекурсивно вложены в своих родителей.
func BuildCommentTree(comments []*models.Comment) []models.Comment {
	childrenMap := make(map[int][]models.Comment)
	var roots []models.Comment

	for _, c := range comments {
		if c.ParentId == nil {
			roots = append(roots, *c)
		} else {
			childrenMap[*c.ParentId] = append(childrenMap[*c.ParentId], *c)
		}
	}

	for i := range roots {
		attachChildren(&roots[i], childrenMap)
	}

	if roots == nil {
		return []models.Comment{}
	}
	return roots
}

// attachChildren рекурсивно прикрепляет дочерние комментарии к родителю.
func attachChildren(parent *models.Comment, childrenMap map[int][]models.Comment) {
	children, ok := childrenMap[parent.Id]
	if !ok {
		return
	}
	parent.Replies = children
	for i := range parent.Replies {
		attachChildren(&parent.Replies[i], childrenMap)
	}
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

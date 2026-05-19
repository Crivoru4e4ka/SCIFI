package services

import (
	"project-MVP/db"
	"project-MVP/models"
)

// CommentStore инкапсулирует операции с комментариями.
type CommentStore struct {
	DB db.DBPool
}

// NewCommentStore создает новый экземпляр CommentStore.
func NewCommentStore(database db.DBPool) *CommentStore {
	return &CommentStore{DB: database}
}

// AddComment создает новый комментарий.
func (s *CommentStore) AddComment(c *models.Comment) (int, error) {
	var id int
	query := `
		INSERT INTO comments (entity_id, entity_type, user_id, parent_id, content)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err := s.DB.QueryRow(query, c.EntityId, c.EntityType, c.UserId, c.ParentId, c.Content).Scan(&id)
	return id, err
}

// GetCommentsByEntity возвращает комментарии сущности в виде дерева.
func (s *CommentStore) GetCommentsByEntity(entityType string, entityId int) ([]models.Comment, error) {
	query := `
		SELECT c.id, c.entity_type, c.entity_id, c.user_id, u.full_name, 
		       c.parent_id, COALESCE(c.content, '') as content, 
               c.created_at, c.updated_at, c.deleted_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.entity_type = $1 AND c.entity_id = $2
		ORDER BY c.created_at ASC`

	rows, err := s.DB.Query(query, entityType, entityId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	allComments := []*models.Comment{}
	for rows.Next() {
		c := &models.Comment{}
		err := rows.Scan(&c.Id, &c.EntityType, &c.EntityId, &c.UserId, &c.UserName,
			&c.ParentId, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt)
		if err != nil {
			return nil, err
		}

		if c.DeletedAt != nil {
			c.Content = "Комментарий удален"
		}

		c.Replies = []models.Comment{}
		allComments = append(allComments, c)
	}

	return BuildCommentTree(allComments), nil
}

// SoftDeleteComment мягко удаляет комментарий.
func (s *CommentStore) SoftDeleteComment(id int) error {
	query := `UPDATE comments SET deleted_at = NOW(), content = NULL WHERE id = $1`
	_, err := s.DB.Exec(query, id)
	return err
}

// GetCommentRaw возвращает сырые данные комментария.
func (s *CommentStore) GetCommentRaw(id int) (*models.Comment, error) {
	var c models.Comment
	err := s.DB.QueryRow("SELECT id, user_id, entity_id, entity_type FROM comments WHERE id = $1", id).
		Scan(&c.Id, &c.UserId, &c.EntityId, &c.EntityType)
	return &c, err
}

// UpdateComment обновляет содержимое комментария.
func (s *CommentStore) UpdateComment(id int, content string) error {
	_, err := s.DB.Exec("UPDATE comments SET content = $1, updated_at = NOW() WHERE id = $2", content, id)
	return err
}

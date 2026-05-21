package services

import (
	"project-MVP/models"
)


// AddComment обёртка над DefaultCommentStore.
func AddComment(c *models.Comment) (int, error) {
	return DefaultCommentStore.AddComment(c)
}

// GetCommentsByEntity обёртка над DefaultCommentStore.
func GetCommentsByEntity(entityType string, entityId int) ([]models.Comment, error) {
	return DefaultCommentStore.GetCommentsByEntity(entityType, entityId)
}

// SoftDeleteComment обёртка над DefaultCommentStore.
func SoftDeleteComment(id int) error {
	return DefaultCommentStore.SoftDeleteComment(id)
}

// GetCommentRaw обёртка над DefaultCommentStore.
func GetCommentRaw(id int) (*models.Comment, error) {
	return DefaultCommentStore.GetCommentRaw(id)
}

// UpdateComment обёртка над DefaultCommentStore.
func UpdateComment(id int, content string) error {
	return DefaultCommentStore.UpdateComment(id, content)
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

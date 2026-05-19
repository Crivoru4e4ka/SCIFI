package services

import (
	"testing"
	"time"

	"project-MVP/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBuildCommentTree_Empty проверяет, что пустой список
// комментариев возвращает пустой слайс (не nil).
func TestBuildCommentTree_Empty(t *testing.T) {
	result := BuildCommentTree([]*models.Comment{})
	assert.Empty(t, result)
	assert.NotNil(t, result)
}

// TestBuildCommentTree_FlatList проверяет обработку списка
// корневых комментариев без ответов.
func TestBuildCommentTree_FlatList(t *testing.T) {
	comments := []*models.Comment{
		{Id: 1, Content: "First", ParentId: nil, Replies: []models.Comment{}},
		{Id: 2, Content: "Second", ParentId: nil, Replies: []models.Comment{}},
	}

	result := BuildCommentTree(comments)

	assert.Len(t, result, 2)
	assert.Equal(t, "First", result[0].Content)
	assert.Equal(t, "Second", result[1].Content)
	assert.Empty(t, result[0].Replies)
	assert.Empty(t, result[1].Replies)
}

// TestBuildCommentTree_NestedReplies проверяет корректное
// построение дерева с вложенными ответами.
func TestBuildCommentTree_NestedReplies(t *testing.T) {
	parentID := 1
	replyID := 2
	comments := []*models.Comment{
		{Id: 1, Content: "Parent", ParentId: nil, Replies: []models.Comment{}},
		{Id: 2, Content: "Reply", ParentId: &replyID, Replies: []models.Comment{}},
	}
	// Исправляем: reply должен ссылаться на parent
	comments[1].ParentId = &parentID

	result := BuildCommentTree(comments)

	assert.Len(t, result, 1)
	assert.Equal(t, "Parent", result[0].Content)
	assert.Len(t, result[0].Replies, 1)
	assert.Equal(t, "Reply", result[0].Replies[0].Content)
}

// TestBuildCommentTree_MultipleReplies проверяет, что несколько
// ответов к одному родителю группируются корректно.
func TestBuildCommentTree_MultipleReplies(t *testing.T) {
	parentID := 1
	comments := []*models.Comment{
		{Id: 1, Content: "Parent", ParentId: nil, Replies: []models.Comment{}},
		{Id: 2, Content: "Reply A", ParentId: &parentID, Replies: []models.Comment{}},
		{Id: 3, Content: "Reply B", ParentId: &parentID, Replies: []models.Comment{}},
	}

	result := BuildCommentTree(comments)

	assert.Len(t, result, 1)
	assert.Len(t, result[0].Replies, 2)
	assert.Equal(t, "Reply A", result[0].Replies[0].Content)
	assert.Equal(t, "Reply B", result[0].Replies[1].Content)
}

// TestBuildCommentTree_DeepNesting проверяет поддержку
// многоуровневой вложенности (reply на reply).
func TestBuildCommentTree_DeepNesting(t *testing.T) {
	rootID := 1
	reply1ID := 2
	comments := []*models.Comment{
		{Id: 1, Content: "Root", ParentId: nil, Replies: []models.Comment{}},
		{Id: 2, Content: "Level 1", ParentId: &rootID, Replies: []models.Comment{}},
		{Id: 3, Content: "Level 2", ParentId: &reply1ID, Replies: []models.Comment{}},
	}

	result := BuildCommentTree(comments)

	require.Len(t, result, 1)
	assert.Len(t, result[0].Replies, 1)
	assert.Len(t, result[0].Replies[0].Replies, 1)
	assert.Equal(t, "Level 2", result[0].Replies[0].Replies[0].Content)
}

// TestBuildCommentTree_OrphanReply проверяет, что ответ
// с несуществующим родителем игнорируется (не попадает в результат).
func TestBuildCommentTree_OrphanReply(t *testing.T) {
	missingParentID := 999
	comments := []*models.Comment{
		{Id: 1, Content: "Root", ParentId: nil, Replies: []models.Comment{}},
		{Id: 2, Content: "Orphan", ParentId: &missingParentID, Replies: []models.Comment{}},
	}

	result := BuildCommentTree(comments)

	assert.Len(t, result, 1)
	assert.Len(t, result[0].Replies, 0)
}

// TestBuildCommentTree_DeletedContent проверяет, что логика
// построения дерева сохраняет модификацию контента для удаленных
// комментариев (установленную до вызова BuildCommentTree).
func TestBuildCommentTree_DeletedContent(t *testing.T) {
	now := time.Now()
	comments := []*models.Comment{
		{
			Id:        1,
			Content:   "Комментарий удален",
			ParentId:  nil,
			DeletedAt: &now,
			Replies:   []models.Comment{},
		},
	}

	result := BuildCommentTree(comments)

	assert.Len(t, result, 1)
	assert.Equal(t, "Комментарий удален", result[0].Content)
}

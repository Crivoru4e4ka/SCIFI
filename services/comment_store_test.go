package services

import (
	"testing"
	"time"

	"project-MVP/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCommentStore_AddComment_Success проверяет создание комментария.
func TestCommentStore_AddComment_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewCommentStore(db)

	mock.ExpectQuery(`INSERT INTO comments \(entity_id, entity_type, user_id, parent_id, content\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING id`).
		WithArgs(1, "task", 1, sqlmock.AnyArg(), "Hello").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(5))

	id, err := store.AddComment(&models.Comment{
		EntityId:   1,
		EntityType: "task",
		UserId:     1,
		Content:    "Hello",
	})

	require.NoError(t, err)
	assert.Equal(t, 5, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestCommentStore_GetCommentsByEntity_Success проверяет получение
// дерева комментариев для сущности.
func TestCommentStore_GetCommentsByEntity_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewCommentStore(db)
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "entity_type", "entity_id", "user_id", "full_name",
		"parent_id", "content", "created_at", "updated_at", "deleted_at",
	}).
		AddRow(1, "task", 1, 1, "Alice", nil, "Root", createdAt, nil, nil).
		AddRow(2, "task", 1, 2, "Bob", 1, "Reply", createdAt, nil, nil)

	mock.ExpectQuery(`SELECT c.id, c.entity_type, c.entity_id, c.user_id, u.full_name, c.parent_id, COALESCE\(c.content, ''\) as content, c.created_at, c.updated_at, c.deleted_at FROM comments c JOIN users u ON c.user_id = u.id WHERE c.entity_type = \$1 AND c.entity_id = \$2 ORDER BY c.created_at ASC`).
		WithArgs("task", 1).
		WillReturnRows(rows)

	comments, err := store.GetCommentsByEntity("task", 1)

	require.NoError(t, err)
	assert.Len(t, comments, 1)
	assert.Equal(t, "Root", comments[0].Content)
	assert.Len(t, comments[0].Replies, 1)
	assert.Equal(t, "Reply", comments[0].Replies[0].Content)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestCommentStore_GetCommentsByEntity_Deleted проверяет, что удаленные
// комментарии отображаются с замещающим текстом.
func TestCommentStore_GetCommentsByEntity_Deleted(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewCommentStore(db)
	createdAt := time.Now()
	deletedAt := createdAt.Add(time.Hour)

	rows := sqlmock.NewRows([]string{
		"id", "entity_type", "entity_id", "user_id", "full_name",
		"parent_id", "content", "created_at", "updated_at", "deleted_at",
	}).AddRow(1, "task", 1, 1, "Alice", nil, "Original", createdAt, nil, deletedAt)

	mock.ExpectQuery(`SELECT c.id, c.entity_type, c.entity_id, c.user_id, u.full_name, c.parent_id, COALESCE\(c.content, ''\) as content, c.created_at, c.updated_at, c.deleted_at FROM comments c JOIN users u ON c.user_id = u.id WHERE c.entity_type = \$1 AND c.entity_id = \$2 ORDER BY c.created_at ASC`).
		WithArgs("task", 1).
		WillReturnRows(rows)

	comments, err := store.GetCommentsByEntity("task", 1)

	require.NoError(t, err)
	assert.Len(t, comments, 1)
	assert.Equal(t, "Комментарий удален", comments[0].Content)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestCommentStore_SoftDeleteComment_Success проверяет мягкое удаление.
func TestCommentStore_SoftDeleteComment_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewCommentStore(db)

	mock.ExpectExec(`UPDATE comments SET deleted_at = NOW\(\), content = NULL WHERE id = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = store.SoftDeleteComment(1)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestCommentStore_GetCommentRaw_Success проверяет получение сырых данных комментария.
func TestCommentStore_GetCommentRaw_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewCommentStore(db)

	mock.ExpectQuery(`SELECT id, user_id, entity_id, entity_type FROM comments WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "entity_id", "entity_type"}).AddRow(1, 2, 3, "task"))

	c, err := store.GetCommentRaw(1)

	require.NoError(t, err)
	assert.Equal(t, 1, c.Id)
	assert.Equal(t, 2, c.UserId)
	assert.Equal(t, 3, c.EntityId)
	assert.Equal(t, "task", c.EntityType)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestCommentStore_UpdateComment_Success проверяет обновление комментария.
func TestCommentStore_UpdateComment_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewCommentStore(db)

	mock.ExpectExec(`UPDATE comments SET content = \$1, updated_at = NOW\(\) WHERE id = \$2`).
		WithArgs("Updated", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = store.UpdateComment(1, "Updated")

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

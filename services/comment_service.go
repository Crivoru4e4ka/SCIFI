package services

import (
	"project-MVP/db"
	"project-MVP/models"
)

func CreateComment(comment models.Comment) (models.Comment, error) {
	if _, err := GetTaskByID(comment.TaskId); err != nil {
		return models.Comment{}, err
	}
	if _, err := GetUserByID(comment.UserId); err != nil {
		return models.Comment{}, err
	}

	query := `INSERT INTO comments (task_id, user_id, content, created_at) VALUES ($1,$2,$3,$4) RETURNING id`
	if err := db.DB.QueryRow(query, comment.TaskId, comment.UserId, comment.Content, comment.CreatedAt).Scan(&comment.Id); err != nil {
		return models.Comment{}, err
	}

	return comment, nil
}

func GetCommentsByTask(taskId int) ([]models.Comment, error) {
	if _, err := GetTaskByID(taskId); err != nil {
		return nil, err
	}

	rows, err := db.DB.Query(`SELECT id, task_id, user_id, content, created_at FROM comments WHERE task_id=$1 ORDER BY created_at`, taskId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.Id, &c.TaskId, &c.UserId, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

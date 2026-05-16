package services

import (
	"errors"
	"strings"

	"project-MVP/db"
	"project-MVP/models"
)

var allowedProjectMemberRoles = map[string]bool{
	"lead":       true, // Руководитель проекта (PI)
	"supervisor": true, // Научный руководитель
	"researcher": true, // Исследователь (основной исполнитель)
	"analyst":    true, // Аналитик данных
	"reviewer":   true, // Рецензент (проверяющий)
	"watcher":    true, // Наблюдатель
}

func isValidProjectMemberRole(role string) bool {
	return allowedProjectMemberRoles[strings.ToLower(strings.TrimSpace(role))]
}

func AddProjectMember(member models.ProjectMember) (models.ProjectMember, error) {
	member.Role = strings.ToLower(strings.TrimSpace(member.Role))
	if !isValidProjectMemberRole(member.Role) {
		return models.ProjectMember{}, errors.New("invalid project member role")
	}

	if _, err := GetProjectByID(member.ProjectId); err != nil {
		return models.ProjectMember{}, err
	}
	if _, err := GetUserByID(member.UserId); err != nil {
		return models.ProjectMember{}, err
	}

	if exists, err := IsUserInProject(member.ProjectId, member.UserId); err != nil {
		return models.ProjectMember{}, err
	} else if exists {
		return models.ProjectMember{}, errors.New("user is already a project member")
	}

	query := `INSERT INTO project_members (project_id, user_id, role) VALUES ($1,$2,$3) RETURNING id`
	if err := db.DB.QueryRow(query, member.ProjectId, member.UserId, member.Role).Scan(&member.Id); err != nil {
		return models.ProjectMember{}, err
	}

	return member, nil
}

func GetProjectMembers(projectId int) ([]models.ProjectMember, error) {
	if _, err := GetProjectByID(projectId); err != nil {
		return nil, err
	}

	rows, err := db.DB.Query(`SELECT id, project_id, user_id, role FROM project_members WHERE project_id=$1`, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.ProjectMember
	for rows.Next() {
		var pm models.ProjectMember
		if err := rows.Scan(&pm.Id, &pm.ProjectId, &pm.UserId, &pm.Role); err != nil {
			return nil, err
		}
		members = append(members, pm)
	}

	return members, nil
}

func IsUserInProject(projectId, userId int) (bool, error) {
	var count int
	row := db.DB.QueryRow(`SELECT COUNT(*) FROM project_members WHERE project_id=$1 AND user_id=$2`, projectId, userId)
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

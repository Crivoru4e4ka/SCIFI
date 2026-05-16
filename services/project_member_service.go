package services

import (
	"database/sql"
	"errors"
	"strings"

	"project-MVP/db"
	"project-MVP/models"
)

func isValidProjectMemberRole(role string) bool {
	return IsValidProjectRole(role)
}

// isTeamManagedProject проверяет, что проект управляется командой
func isTeamManagedProject(projectId int) (bool, error) {
	var execType string
	err := db.DB.QueryRow(`SELECT execution_type FROM projects WHERE id = $1`, projectId).Scan(&execType)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, ErrNotFound
		}
		return false, err
	}
	return execType == "team", nil
}

// getProjectRoleIDByName возвращает role_id из таблицы roles по имени проектной роли
func getProjectRoleIDByName(role string) (*int, error) {
	role = strings.ToLower(strings.TrimSpace(role))
	if role == "" {
		return nil, nil
	}
	id, err := getRoleIDByName(role)
	if err != nil {
		return nil, err
	}
	return &id, nil
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

	// Блокируем ручное управление участниками для team-managed проектов
	if isTeam, _ := isTeamManagedProject(member.ProjectId); isTeam {
		return models.ProjectMember{}, errors.New("cannot manually manage members of a team-managed project")
	}

	if exists, err := IsUserInProject(member.ProjectId, member.UserId); err != nil {
		return models.ProjectMember{}, err
	} else if exists {
		return models.ProjectMember{}, errors.New("user is already a project member")
	}

	// Автоматически определяем role_id по строковой роли
	roleID := member.RoleId
	if roleID == nil || *roleID <= 0 {
		mappedID, err := getProjectRoleIDByName(member.Role)
		if err == nil && mappedID != nil {
			roleID = mappedID
		}
	}

	query := `INSERT INTO project_members (project_id, user_id, role, role_id) VALUES ($1,$2,$3,$4) RETURNING id`
	if err := db.DB.QueryRow(query, member.ProjectId, member.UserId, member.Role, roleID).Scan(&member.Id); err != nil {
		return models.ProjectMember{}, err
	}

	return member, nil
}

func RemoveProjectMember(projectId, userId int) error {
	if _, err := GetProjectByID(projectId); err != nil {
		return err
	}

	// Блокируем ручное управление участниками для team-managed проектов
	if isTeam, _ := isTeamManagedProject(projectId); isTeam {
		return errors.New("cannot manually manage members of a team-managed project")
	}

	res, err := db.DB.Exec(`DELETE FROM project_members WHERE project_id=$1 AND user_id=$2`, projectId, userId)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("member not found")
	}
	return nil
}

func UpdateProjectMemberRole(projectId, userId int, role string, roleID *int) error {
	role = strings.ToLower(strings.TrimSpace(role))
	if !isValidProjectMemberRole(role) {
		return errors.New("invalid project member role")
	}

	if _, err := GetProjectByID(projectId); err != nil {
		return err
	}

	// Блокируем ручное управление участниками для team-managed проектов
	if isTeam, _ := isTeamManagedProject(projectId); isTeam {
		return errors.New("cannot manually manage members of a team-managed project")
	}

	// Автоматически определяем role_id если не передан
	if roleID == nil || *roleID <= 0 {
		mappedID, err := getProjectRoleIDByName(role)
		if err == nil && mappedID != nil {
			roleID = mappedID
		}
	}

	res, err := db.DB.Exec(`UPDATE project_members SET role=$1, role_id=$2 WHERE project_id=$3 AND user_id=$4`, role, roleID, projectId, userId)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("member not found")
	}
	return nil
}

func GetProjectMembers(projectId int) ([]models.ProjectMember, error) {
	if _, err := GetProjectByID(projectId); err != nil {
		return nil, err
	}

	rows, err := db.DB.Query(`SELECT id, project_id, user_id, role, role_id FROM project_members WHERE project_id=$1`, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.ProjectMember
	for rows.Next() {
		var pm models.ProjectMember
		var roleID sql.NullInt64
		if err := rows.Scan(&pm.Id, &pm.ProjectId, &pm.UserId, &pm.Role, &roleID); err != nil {
			return nil, err
		}
		if roleID.Valid {
			v := int(roleID.Int64)
			pm.RoleId = &v
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

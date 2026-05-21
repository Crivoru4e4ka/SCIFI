package services

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"project-MVP/db"
	"project-MVP/models"
)


// ValidateAndNormalizeProject проверяет и нормализует поля проекта.
// Вынесена в отдельную функцию для unit-тестирования без обращения к БД.
func ValidateAndNormalizeProject(project *models.Project) error {
	project.Name = strings.TrimSpace(project.Name)
	project.Description = strings.TrimSpace(project.Description)
	project.Status = strings.ToLower(strings.TrimSpace(project.Status))
	project.ExecutionType = strings.ToLower(strings.TrimSpace(project.ExecutionType))

	if project.Name == "" {
		return errors.New("project name is required")
	}
	if project.CreatedBy <= 0 {
		return errors.New("created_by is required")
	}
	if project.StartDate.IsZero() {
		project.StartDate = time.Now()
	}
	if project.Status == "" {
		project.Status = "active"
	}
	if project.Visibility == "" {
		project.Visibility = "closed"
	}
	if project.ExecutionType == "" {
		project.ExecutionType = "manual"
	}

	if project.ExecutionType != "manual" && project.ExecutionType != "team" {
		return errors.New("invalid execution type")
	}

	if project.ExecutionType == "team" && (project.TeamId == nil || *project.TeamId <= 0) {
		return errors.New("team_id is required for team execution type")
	}

	return nil
}

// GetProjects обёртка над DefaultProjectStore.
func GetProjects() ([]models.Project, error) {
	return DefaultProjectStore.GetProjects()
}

// GetUserProjects обёртка над DefaultProjectStore.
func GetUserProjects(userId int) ([]models.Project, error) {
	return DefaultProjectStore.GetUserProjects(userId)
}

// GetProjectByID обёртка над DefaultProjectStore.
func GetProjectByID(id int) (models.Project, error) {
	return DefaultProjectStore.GetProjectByID(id)
}

// GetProjectByName обёртка над DefaultProjectStore.
func GetProjectByName(name string) (models.Project, error) {
	return DefaultProjectStore.GetProjectByName(name)
}

// GetProjectProgress обёртка над DefaultProjectStore.
func GetProjectProgress(projectId int) (float64, error) {
	return DefaultProjectStore.GetProjectProgress(projectId)
}

// GetProjectAssignableUsers обёртка над DefaultProjectStore.
func GetProjectAssignableUsers(projectID int) ([]models.User, error) {
	return DefaultProjectStore.GetProjectAssignableUsers(projectID)
}

// UpdateProject обёртка над DefaultProjectStore.
func UpdateProject(id int, p models.Project) error {
	return DefaultProjectStore.UpdateProject(id, p)
}

// DeleteProject обёртка над DefaultProjectStore.
func DeleteProject(id int) error {
	return DefaultProjectStore.DeleteProject(id)
}

// CreateProject обёртка над DefaultProjectStore.
func CreateProject(project models.Project) (models.Project, error) {
	return DefaultProjectStore.CreateProject(project)
}

// getTeamMembersTx получает участников команды внутри транзакции.
func getTeamMembersTx(tx *sql.Tx, teamID int) ([]models.TeamMemberInfo, error) {
	rows, err := tx.Query(`
		SELECT tm.user_id, u.full_name, u.email, tm.role
		FROM team_members tm
		JOIN users u ON u.id = tm.user_id
		WHERE tm.team_id = $1`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.TeamMemberInfo
	for rows.Next() {
		var m models.TeamMemberInfo
		if err := rows.Scan(&m.UserID, &m.FullName, &m.Email, &m.Role); err != nil {
			continue
		}
		members = append(members, m)
	}
	return members, nil
}

// resolveTeamMemberRole преобразует роль из team_members в валидную проектную роль с role_id.
// Если роль не является валидной проектной ролью — возвращает researcher.
func resolveTeamMemberRole(teamRole string) (string, interface{}) {
	roleName := strings.ToLower(strings.TrimSpace(teamRole))
	if !IsValidProjectRole(roleName) {
		roleName = "researcher"
	}
	roleID, err := getRoleIDByName(roleName)
	if err != nil || roleID <= 0 {
		return roleName, nil
	}
	return roleName, roleID
}

// SyncTeamMembersToProject синхронизирует участников команды с проектом
// (вызывать при добавлении нового участника в команду, если проект уже существует).
func SyncTeamMembersToProject(projectID, teamID int) error {
	project, err := GetProjectByID(projectID)
	if err != nil {
		return err
	}
	if project.ExecutionType != "team" || project.TeamId == nil || *project.TeamId != teamID {
		return errors.New("project is not linked to this team")
	}

	teamMembers, err := GetTeamMembers(teamID)
	if err != nil {
		return err
	}

	for _, tm := range teamMembers {
		exists, _ := IsUserInProject(projectID, tm.UserID)
		if exists {
			continue
		}
		memberRole, memberRoleID := resolveTeamMemberRole(tm.Role)
		_, err := db.DB.Exec(`INSERT INTO project_members (project_id, user_id, role, role_id)
			VALUES ($1,$2,$3,$4)
			ON CONFLICT (project_id, user_id) DO NOTHING`,
			projectID, tm.UserID, memberRole, memberRoleID)
		if err != nil {
			return err
		}
	}
	return nil
}

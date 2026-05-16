package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"project-MVP/db"
	"project-MVP/models"
)

var (
	// ErrNoPermission возвращается при отсутствии права
	ErrNoPermission = errors.New("insufficient permissions")
)

// GetUserSystemRole возвращает системную роль пользователя
func GetUserSystemRole(userID int) (string, error) {
	var role string
	err := db.DB.QueryRow(`SELECT role FROM users WHERE id = $1`, userID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrNotFound
		}
		return "", err
	}
	return role, nil
}

// GetUserProjectRoleID возвращает ID проектной роли пользователя
func GetUserProjectRoleID(userID, projectID int) (*int, error) {
	var roleID sql.NullInt64
	err := db.DB.QueryRow(`SELECT role_id FROM project_members WHERE user_id = $1 AND project_id = $2`, userID, projectID).Scan(&roleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Пользователь не участник проекта
		}
		return nil, err
	}
	if roleID.Valid {
		v := int(roleID.Int64)
		return &v, nil
	}
	return nil, nil
}

// GetRolePermissions возвращает список кодов прав для роли
func GetRolePermissions(roleID int) ([]string, error) {
	rows, err := db.DB.Query(`
		SELECT p.code FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = $1
	`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			continue
		}
		perms = append(perms, code)
	}
	return perms, nil
}

// GetUserPermissions возвращает все права пользователя (системные + проектные)
func GetUserPermissions(userID, projectID int) ([]string, error) {
	// Проверяем системную роль
	sysRole, err := GetUserSystemRole(userID)
	if err != nil {
		return nil, err
	}

	// Admin имеет все права
	if sysRole == "admin" {
		return getAllPermissionCodes(), nil
	}

	permsMap := make(map[string]bool)

	// Добавляем права системной роли
	if sysRole != "" {
		roleID, err := getRoleIDByName(sysRole)
		if err == nil && roleID > 0 {
			rolePerms, _ := GetRolePermissions(roleID)
			for _, p := range rolePerms {
				permsMap[p] = true
			}
		}
	}

	// Добавляем права проектной роли
	if projectID > 0 {
		projRoleID, err := GetUserProjectRoleID(userID, projectID)
		if err == nil && projRoleID != nil {
			rolePerms, _ := GetRolePermissions(*projRoleID)
			for _, p := range rolePerms {
				permsMap[p] = true
			}
		}
	}

	var result []string
	for p := range permsMap {
		result = append(result, p)
	}
	return result, nil
}

// HasPermission проверяет, есть ли у пользователя право
func HasPermission(userID, projectID int, permissionCode string) bool {
	perms, err := GetUserPermissions(userID, projectID)
	if err != nil {
		log.Printf("HasPermission error: %v", err)
		return false
	}
	for _, p := range perms {
		if p == permissionCode {
			return true
		}
	}
	return false
}

// CheckPermission возвращает ошибку если права нет
func CheckPermission(userID, projectID int, permissionCode string) error {
	if !HasPermission(userID, projectID, permissionCode) {
		return ErrNoPermission
	}
	return nil
}

// AssignProjectRole назначает роль участнику проекта
func AssignProjectRole(projectID, userID, roleID int) error {
	// Блокируем ручное управление участниками для team-managed проектов
	var execType string
	err := db.DB.QueryRow(`SELECT execution_type FROM projects WHERE id = $1`, projectID).Scan(&execType)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	if execType == "team" {
		return errors.New("cannot manually manage members of a team-managed project")
	}

	// Проверяем существование участника
	var exists bool
	err = db.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM project_members WHERE project_id = $1 AND user_id = $2)`, projectID, userID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user is not a member of this project")
	}

	// Проверяем существование роли и получаем её имя
	var roleName string
	err = db.DB.QueryRow(`SELECT name FROM roles WHERE id = $1`, roleID).Scan(&roleName)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("role not found")
		}
		return err
	}

	_, err = db.DB.Exec(`UPDATE project_members SET role_id = $1, role = $2 WHERE project_id = $3 AND user_id = $4`, roleID, roleName, projectID, userID)
	return err
}

// GetProjectMembersWithRoles возвращает участников проекта с ролями
func GetProjectMembersWithRoles(projectID int) ([]models.ProjectMemberWithRole, error) {
	query := `
		SELECT pm.id, pm.user_id, pm.project_id, pm.role_id,
			COALESCE(r.name, pm.role, '') as role_name,
			COALESCE(r.description, '') as role_desc,
			u.full_name, u.email
		FROM project_members pm
		JOIN users u ON u.id = pm.user_id
		LEFT JOIN roles r ON r.id = pm.role_id
		WHERE pm.project_id = $1
		ORDER BY pm.id`

	rows, err := db.DB.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.ProjectMemberWithRole
	for rows.Next() {
		var m models.ProjectMemberWithRole
		var roleID sql.NullInt64
		err := rows.Scan(&m.Id, &m.UserId, &m.ProjectId, &roleID, &m.RoleName, &m.RoleDesc, &m.UserName, &m.UserEmail)
		if err != nil {
			continue
		}
		if roleID.Valid {
			v := int(roleID.Int64)
			m.RoleId = &v
		}
		list = append(list, m)
	}
	return list, nil
}

// GetAllRoles возвращает все роли
func GetAllRoles() ([]models.Role, error) {
	rows, err := db.DB.Query(`SELECT id, name, description, is_system, created_at FROM roles ORDER BY is_system DESC, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var r models.Role
		if err := rows.Scan(&r.Id, &r.Name, &r.Description, &r.IsSystem, &r.CreatedAt); err != nil {
			continue
		}
		roles = append(roles, r)
	}
	return roles, nil
}

// GetAllPermissions возвращает все права
func GetAllPermissions() ([]models.Permission, error) {
	rows, err := db.DB.Query(`SELECT id, code, name, category FROM permissions ORDER BY category, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []models.Permission
	for rows.Next() {
		var p models.Permission
		if err := rows.Scan(&p.Id, &p.Code, &p.Name, &p.Category); err != nil {
			continue
		}
		perms = append(perms, p)
	}
	return perms, nil
}

// LogAudit записывает событие в аудит-лог
func LogAudit(userID, projectID int, action, entityType string, entityID int, details string) {
	query := `INSERT INTO audit_log (user_id, project_id, action, entity_type, entity_id, details) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := db.DB.Exec(query, nullInt(userID), nullInt(projectID), action, entityType, nullInt(entityID), details)
	if err != nil {
		log.Printf("Audit log error: %v", err)
	}
}

// GetAuditLog возвращает аудит-лог проекта
func GetAuditLog(projectID int, limit int) ([]models.AuditLog, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `
		SELECT al.id, al.user_id, al.project_id, al.action, al.entity_type, al.entity_id, al.details, al.created_at,
			u.full_name
		FROM audit_log al
		LEFT JOIN users u ON u.id = al.user_id
		WHERE al.project_id = $1
		ORDER BY al.created_at DESC
		LIMIT $2`

	rows, err := db.DB.Query(query, projectID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var l models.AuditLog
		var userID sql.NullInt64
		var entityID sql.NullInt64
		var userName sql.NullString
		err := rows.Scan(&l.Id, &userID, &l.ProjectId, &l.Action, &l.EntityType, &entityID, &l.Details, &l.CreatedAt, &userName)
		if err != nil {
			continue
		}
		if userID.Valid {
			v := int(userID.Int64)
			l.UserId = &v
		}
		if entityID.Valid {
			v := int(entityID.Int64)
			l.EntityId = &v
		}
		if userName.Valid && userName.String != "" {
			l.Details = fmt.Sprintf("%s (%s)", l.Details, userName.String)
		}
		logs = append(logs, l)
	}
	return logs, nil
}

// Вспомогательные функции

func getAllPermissionCodes() []string {
	rows, err := db.DB.Query(`SELECT code FROM permissions`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var codes []string
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err == nil {
			codes = append(codes, c)
		}
	}
	return codes
}

func getRoleIDByName(name string) (int, error) {
	var id int
	err := db.DB.QueryRow(`SELECT id FROM roles WHERE name = $1`, name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

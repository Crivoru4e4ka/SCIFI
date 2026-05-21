package services

import (
	"fmt"

	"project-MVP/db"
	"project-MVP/models"
)

var (
	// ErrNoPermission возвращается при отсутствии права
	ErrNoPermission = fmt.Errorf("insufficient permissions")
)


// GetUserSystemRole обёртка над DefaultRBACStore.
func GetUserSystemRole(userID int) (string, error) {
	return DefaultRBACStore.GetUserSystemRole(userID)
}

// GetUserProjectRoleID обёртка над DefaultRBACStore.
func GetUserProjectRoleID(userID, projectID int) (*int, error) {
	return DefaultRBACStore.GetUserProjectRoleID(userID, projectID)
}

// GetRolePermissions обёртка над DefaultRBACStore.
func GetRolePermissions(roleID int) ([]string, error) {
	return DefaultRBACStore.GetRolePermissions(roleID)
}

// GetUserPermissions обёртка над DefaultRBACStore.
func GetUserPermissions(userID, projectID int) ([]string, error) {
	return DefaultRBACStore.GetUserPermissions(userID, projectID)
}

// HasPermission обёртка над DefaultRBACStore.
func HasPermission(userID, projectID int, permissionCode string) bool {
	return DefaultRBACStore.HasPermission(userID, projectID, permissionCode)
}

// CheckPermission обёртка над DefaultRBACStore.
func CheckPermission(userID, projectID int, permissionCode string) error {
	return DefaultRBACStore.CheckPermission(userID, projectID, permissionCode)
}

// AssignProjectRole обёртка над DefaultRBACStore.
func AssignProjectRole(projectID, userID, roleID int) error {
	return DefaultRBACStore.AssignProjectRole(projectID, userID, roleID)
}

// GetProjectMembersWithRoles обёртка над DefaultRBACStore.
func GetProjectMembersWithRoles(projectID int) ([]models.ProjectMemberWithRole, error) {
	return DefaultRBACStore.GetProjectMembersWithRoles(projectID)
}

// GetAllRoles обёртка над DefaultRBACStore.
func GetAllRoles() ([]models.Role, error) {
	return DefaultRBACStore.GetAllRoles()
}

// GetAllPermissions обёртка над DefaultRBACStore.
func GetAllPermissions() ([]models.Permission, error) {
	return DefaultRBACStore.GetAllPermissions()
}

// LogAudit обёртка над DefaultRBACStore.
func LogAudit(userID, projectID int, action, entityType string, entityID int, details string) {
	DefaultRBACStore.LogAudit(userID, projectID, action, entityType, entityID, details)
}

// GetAuditLog обёртка над DefaultRBACStore.
func GetAuditLog(projectID int, limit int) ([]models.AuditLog, error) {
	return DefaultRBACStore.GetAuditLog(projectID, limit)
}

// IsAdmin обёртка над DefaultRBACStore.
func IsAdmin(userID int) bool {
	return DefaultRBACStore.IsAdmin(userID)
}

// IsProjectMember обёртка над DefaultRBACStore.
func IsProjectMember(userID, projectID int) bool {
	return DefaultRBACStore.IsProjectMember(userID, projectID)
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

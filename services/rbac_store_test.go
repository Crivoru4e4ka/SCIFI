package services

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRBACStore_GetUserSystemRole_Success проверяет получение системной роли.
func TestRBACStore_GetUserSystemRole_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewRBACStore(db)

	mock.ExpectQuery(`SELECT role FROM users WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("admin"))

	role, err := store.GetUserSystemRole(1)

	require.NoError(t, err)
	assert.Equal(t, "admin", role)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRBACStore_GetUserSystemRole_NotFound проверяет обработку отсутствующего пользователя.
func TestRBACStore_GetUserSystemRole_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewRBACStore(db)

	mock.ExpectQuery(`SELECT role FROM users WHERE id = \$1`).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	_, err = store.GetUserSystemRole(999)

	assert.ErrorIs(t, err, ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRBACStore_IsAdmin_True проверяет определение администратора.
func TestRBACStore_IsAdmin_True(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewRBACStore(db)

	mock.ExpectQuery(`SELECT role FROM users WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("admin"))

	assert.True(t, store.IsAdmin(1))
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRBACStore_IsAdmin_False проверяет определение не-администратора.
func TestRBACStore_IsAdmin_False(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewRBACStore(db)

	mock.ExpectQuery(`SELECT role FROM users WHERE id = \$1`).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("user"))

	assert.False(t, store.IsAdmin(2))
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRBACStore_IsProjectMember_True проверяет членство в проекте.
func TestRBACStore_IsProjectMember_True(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewRBACStore(db)

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM project_members WHERE user_id = \$1 AND project_id = \$2\)`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	assert.True(t, store.IsProjectMember(1, 1))
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRBACStore_IsProjectMember_False проверяет отсутствие членства в проекте.
func TestRBACStore_IsProjectMember_False(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewRBACStore(db)

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM project_members WHERE user_id = \$1 AND project_id = \$2\)`).
		WithArgs(99, 1).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	assert.False(t, store.IsProjectMember(99, 1))
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRBACStore_GetRolePermissions_Success проверяет получение прав роли.
func TestRBACStore_GetRolePermissions_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewRBACStore(db)

	rows := sqlmock.NewRows([]string{"code"}).
		AddRow("project.view").
		AddRow("task.create")

	mock.ExpectQuery(`SELECT p.code FROM permissions p JOIN role_permissions rp ON rp.permission_id = p.id WHERE rp.role_id = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	perms, err := store.GetRolePermissions(1)

	require.NoError(t, err)
	assert.Len(t, perms, 2)
	assert.Contains(t, perms, "project.view")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRBACStore_GetAllRoles_Success проверяет получение всех ролей.
func TestRBACStore_GetAllRoles_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewRBACStore(db)
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "name", "description", "is_system", "created_at"}).
		AddRow(1, "admin", "Admin", true, createdAt).
		AddRow(2, "researcher", "Researcher", false, createdAt)

	mock.ExpectQuery(`SELECT id, name, description, is_system, created_at FROM roles ORDER BY is_system DESC, name`).
		WillReturnRows(rows)

	roles, err := store.GetAllRoles()

	require.NoError(t, err)
	assert.Len(t, roles, 2)
	assert.True(t, roles[0].IsSystem)
	assert.False(t, roles[1].IsSystem)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRBACStore_GetAllPermissions_Success проверяет получение всех прав.
func TestRBACStore_GetAllPermissions_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewRBACStore(db)

	rows := sqlmock.NewRows([]string{"id", "code", "name", "category"}).
		AddRow(1, "project.view", "View", "project").
		AddRow(2, "task.create", "Create", "task")

	mock.ExpectQuery(`SELECT id, code, name, category FROM permissions ORDER BY category, name`).
		WillReturnRows(rows)

	perms, err := store.GetAllPermissions()

	require.NoError(t, err)
	assert.Len(t, perms, 2)
	assert.Equal(t, "project.view", perms[0].Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRBACStore_GetProjectMembersWithRoles_Success проверяет получение
// участников проекта с ролями.
func TestRBACStore_GetProjectMembersWithRoles_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewRBACStore(db)

	rows := sqlmock.NewRows([]string{
		"pm.id", "pm.user_id", "pm.project_id", "pm.role_id",
		"role_name", "role_desc", "u.full_name", "u.email",
	}).AddRow(1, 1, 1, 2, "researcher", "Researcher", "Alice", "a@example.com")

	mock.ExpectQuery(`SELECT pm.id, pm.user_id, pm.project_id, pm.role_id, COALESCE\(r.name, pm.role, ''\) as role_name, COALESCE\(r.description, ''\) as role_desc, u.full_name, u.email FROM project_members pm JOIN users u ON u.id = pm.user_id LEFT JOIN roles r ON r.id = pm.role_id WHERE pm.project_id = \$1 ORDER BY pm.id`).
		WithArgs(1).
		WillReturnRows(rows)

	members, err := store.GetProjectMembersWithRoles(1)

	require.NoError(t, err)
	assert.Len(t, members, 1)
	assert.Equal(t, "Alice", members[0].UserName)
	assert.Equal(t, "researcher", members[0].RoleName)
	require.NotNil(t, members[0].RoleId)
	assert.Equal(t, 2, *members[0].RoleId)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRBACStore_GetAuditLog_Success проверяет получение аудит-лога.
func TestRBACStore_GetAuditLog_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewRBACStore(db)
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{
		"al.id", "al.user_id", "al.project_id", "al.action", "al.entity_type",
		"al.entity_id", "al.details", "al.created_at", "u.full_name",
	}).AddRow(1, 1, 1, "created", "task", 1, "Created task", createdAt, "Alice")

	mock.ExpectQuery(`SELECT al.id, al.user_id, al.project_id, al.action, al.entity_type, al.entity_id, al.details, al.created_at, u.full_name FROM audit_log al LEFT JOIN users u ON u.id = al.user_id WHERE al.project_id = \$1 ORDER BY al.created_at DESC LIMIT \$2`).
		WithArgs(1, 50).
		WillReturnRows(rows)

	logs, err := store.GetAuditLog(1, 0)

	require.NoError(t, err)
	assert.Len(t, logs, 1)
	assert.Equal(t, "created", logs[0].Action)
	assert.Contains(t, logs[0].Details, "Alice")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRBACStore_AssignProjectRole_TeamManaged проверяет блокировку
// ручного назначения ролей для team-managed проектов.
func TestRBACStore_AssignProjectRole_TeamManaged(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewRBACStore(db)

	mock.ExpectQuery(`SELECT execution_type FROM projects WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"execution_type"}).AddRow("team"))

	err = store.AssignProjectRole(1, 2, 3)

	assert.EqualError(t, err, "cannot manually manage members of a team-managed project")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRBACStore_AssignProjectRole_NotMember проверяет отказ при назначении
// роли пользователю, не являющемуся участником проекта.
func TestRBACStore_AssignProjectRole_NotMember(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewRBACStore(db)

	mock.ExpectQuery(`SELECT execution_type FROM projects WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"execution_type"}).AddRow("manual"))

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM project_members WHERE project_id = \$1 AND user_id = \$2\)`).
		WithArgs(1, 99).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	err = store.AssignProjectRole(1, 99, 3)

	assert.EqualError(t, err, "user is not a member of this project")
	assert.NoError(t, mock.ExpectationsWereMet())
}

package services

import (
	"database/sql"
	"testing"

	"project-MVP/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTeamStore_GetTeamByID_Success проверяет получение команды по ID.
func TestTeamStore_GetTeamByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTeamStore(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "created_by"}).
		AddRow(1, "Team Alpha", "Desc", 1)

	mock.ExpectQuery(`SELECT id, name, description, created_by FROM teams WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	team, err := store.GetTeamByID(1)

	require.NoError(t, err)
	assert.Equal(t, 1, team.ID)
	assert.Equal(t, "Team Alpha", team.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTeamStore_GetTeamByID_NotFound проверяет обработку отсутствующей команды.
func TestTeamStore_GetTeamByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTeamStore(db)

	mock.ExpectQuery(`SELECT id, name, description, created_by FROM teams WHERE id = \$1`).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	_, err = store.GetTeamByID(999)

	assert.ErrorIs(t, err, ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTeamStore_GetUserTeams_Success проверяет получение команд пользователя.
func TestTeamStore_GetUserTeams_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTeamStore(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "created_by", "members_count"}).
		AddRow(1, "T1", "D1", 1, 3).
		AddRow(2, "T2", "D2", 2, 5)

	mock.ExpectQuery(`SELECT t.id, t.name, t.description, t.created_by, \(SELECT COUNT\(\*\) FROM team_members tm WHERE tm.team_id = t.id\) as members_count FROM teams t JOIN team_members tm ON t.id = tm.team_id WHERE tm.user_id = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	teams, err := store.GetUserTeams(1)

	require.NoError(t, err)
	assert.Len(t, teams, 2)
	assert.Equal(t, 3, teams[0].MembersCount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTeamStore_CreateTeam_Success проверяет создание команды с транзакцией.
func TestTeamStore_CreateTeam_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTeamStore(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO teams \(name, description, created_by\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
		WithArgs("New Team", "Desc", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))
	mock.ExpectExec(`INSERT INTO team_members \(team_id, user_id, role\) VALUES \(\$1, \$2, 'project_lead'\)`).
		WithArgs(10, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	team, err := store.CreateTeam(models.Team{
		Name:        "New Team",
		Description: "Desc",
		CreatedBy:   1,
	})

	require.NoError(t, err)
	assert.Equal(t, 10, team.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTeamStore_GetTeamMembers_Success проверяет получение участников команды.
func TestTeamStore_GetTeamMembers_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTeamStore(db)

	rows := sqlmock.NewRows([]string{"id", "full_name", "email", "role"}).
		AddRow(1, "Alice", "a@example.com", "project_lead").
		AddRow(2, "Bob", "b@example.com", "researcher")

	mock.ExpectQuery(`SELECT u.id, u.full_name, u.email, tm.role FROM users u JOIN team_members tm ON u.id = tm.user_id WHERE tm.team_id = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	members, err := store.GetTeamMembers(1)

	require.NoError(t, err)
	assert.Len(t, members, 2)
	assert.Equal(t, "Alice", members[0].FullName)
	assert.Equal(t, "researcher", members[1].Role)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTeamStore_RemoveMemberFromTeam_Success проверяет удаление участника.
func TestTeamStore_RemoveMemberFromTeam_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTeamStore(db)

	mock.ExpectExec(`DELETE FROM team_members WHERE team_id = \$1 AND user_id = \$2`).
		WithArgs(1, 2).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = store.RemoveMemberFromTeam(1, 2)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTeamStore_UpdateTeam_Success проверяет обновление команды.
func TestTeamStore_UpdateTeam_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTeamStore(db)

	mock.ExpectExec(`UPDATE teams SET name = \$1, description = \$2 WHERE id = \$3`).
		WithArgs("Updated", "New Desc", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = store.UpdateTeam(1, "Updated", "New Desc")

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTeamStore_AddMemberToTeam_Success проверяет добавление участника по email.
func TestTeamStore_AddMemberToTeam_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTeamStore(db)

	mock.ExpectQuery(`SELECT id, full_name, email FROM users WHERE email = \$1`).
		WithArgs("bob@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "full_name", "email"}).AddRow(2, "Bob", "bob@example.com"))

	mock.ExpectExec(`INSERT INTO team_members \(team_id, user_id, role\) VALUES \(\$1, \$2, 'researcher'\) ON CONFLICT DO NOTHING`).
		WithArgs(1, 2).
		WillReturnResult(sqlmock.NewResult(0, 1))

	member, err := store.AddMemberToTeam(1, "bob@example.com")

	require.NoError(t, err)
	assert.Equal(t, 2, member.UserID)
	assert.Equal(t, "researcher", member.Role)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTeamStore_UpdateMemberRole_InvalidRole проверяет отказ при
// невалидной роли (до обращения к БД).
func TestTeamStore_UpdateMemberRole_InvalidRole(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTeamStore(db)

	err = store.UpdateMemberRole(1, 2, "superuser")
	assert.EqualError(t, err, "invalid team member role: must be a valid project role")
}

// TestTeamStore_DeleteTeam_Success проверяет каскадное удаление команды.
func TestTeamStore_DeleteTeam_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewTeamStore(db)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM team_members WHERE team_id = \$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 3))
	mock.ExpectExec(`DELETE FROM teams WHERE id = \$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = store.DeleteTeam(1)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

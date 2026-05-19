package services

import (
	"database/sql"
	"errors"
	"strings"

	"project-MVP/db"
	"project-MVP/models"
)

// TeamStore инкапсулирует операции с командами.
type TeamStore struct {
	DB db.DBPool
}

// NewTeamStore создает новый экземпляр TeamStore.
func NewTeamStore(database db.DBPool) *TeamStore {
	return &TeamStore{DB: database}
}

// GetTeamByID возвращает команду по ID.
func (s *TeamStore) GetTeamByID(teamID int) (models.Team, error) {
	var t models.Team
	query := `SELECT id, name, description, created_by FROM teams WHERE id = $1`
	err := s.DB.QueryRow(query, teamID).Scan(&t.ID, &t.Name, &t.Description, &t.CreatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return t, ErrNotFound
		}
		return t, err
	}
	return t, nil
}

// GetUserTeams возвращает команды пользователя.
func (s *TeamStore) GetUserTeams(userID int) ([]models.Team, error) {
	query := `
		SELECT t.id, t.name, t.description, t.created_by, 
		(SELECT COUNT(*) FROM team_members tm WHERE tm.team_id = t.id) as members_count
		FROM teams t
		JOIN team_members tm ON t.id = tm.team_id
		WHERE tm.user_id = $1`

	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var t models.Team
		rows.Scan(&t.ID, &t.Name, &t.Description, &t.CreatedBy, &t.MembersCount)
		teams = append(teams, t)
	}
	return teams, nil
}

// CreateTeam создает команду с создателем и приглашенными участниками.
func (s *TeamStore) CreateTeam(t models.Team) (models.Team, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return t, err
	}
	defer tx.Rollback()

	err = tx.QueryRow(`INSERT INTO teams (name, description, created_by) VALUES ($1, $2, $3) RETURNING id`,
		t.Name, t.Description, t.CreatedBy).Scan(&t.ID)
	if err != nil {
		return t, err
	}

	_, err = tx.Exec(`INSERT INTO team_members (team_id, user_id, role) VALUES ($1, $2, 'project_lead')`, t.ID, t.CreatedBy)
	if err != nil {
		return t, err
	}

	for _, email := range t.MemberEmails {
		var invitedID int
		err := tx.QueryRow("SELECT id FROM users WHERE email = $1", strings.TrimSpace(email)).Scan(&invitedID)
		if err != nil {
			continue
		}
		if invitedID == t.CreatedBy {
			continue
		}

		_, err = tx.Exec(`INSERT INTO team_members (team_id, user_id, role) VALUES ($1, $2, 'researcher')`, t.ID, invitedID)
		if err != nil {
			return t, err
		}
	}

	err = tx.Commit()
	return t, err
}

// GetTeamMembers возвращает участников команды.
func (s *TeamStore) GetTeamMembers(teamID int) ([]models.TeamMemberInfo, error) {
	query := `
        SELECT u.id, u.full_name, u.email, tm.role 
        FROM users u
        JOIN team_members tm ON u.id = tm.user_id
        WHERE tm.team_id = $1`

	rows, err := s.DB.Query(query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []models.TeamMemberInfo{}
	for rows.Next() {
		var m models.TeamMemberInfo
		if err := rows.Scan(&m.UserID, &m.FullName, &m.Email, &m.Role); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

// RemoveMemberFromTeam удаляет участника из команды.
func (s *TeamStore) RemoveMemberFromTeam(teamID int, userID int) error {
	_, err := s.DB.Exec("DELETE FROM team_members WHERE team_id = $1 AND user_id = $2", teamID, userID)
	return err
}

// UpdateTeam обновляет информацию о команде.
func (s *TeamStore) UpdateTeam(teamID int, name, description string) error {
	_, err := s.DB.Exec("UPDATE teams SET name = $1, description = $2 WHERE id = $3", name, description, teamID)
	return err
}

// AddMemberToTeam добавляет участника в команду по email.
func (s *TeamStore) AddMemberToTeam(teamID int, email string) (models.TeamMemberInfo, error) {
	var m models.TeamMemberInfo
	var userID int

	err := s.DB.QueryRow("SELECT id, full_name, email FROM users WHERE email = $1", strings.TrimSpace(email)).
		Scan(&userID, &m.FullName, &m.Email)
	if err != nil {
		return m, err
	}

	_, err = s.DB.Exec("INSERT INTO team_members (team_id, user_id, role) VALUES ($1, $2, 'researcher') ON CONFLICT DO NOTHING", teamID, userID)

	m.UserID = userID
	m.Role = "researcher"
	return m, err
}

// UpdateMemberRole обновляет роль участника команды.
func (s *TeamStore) UpdateMemberRole(teamID, userID int, newRole string) error {
	newRole = strings.ToLower(strings.TrimSpace(newRole))
	if !IsValidProjectRole(newRole) {
		return errors.New("invalid team member role: must be a valid project role")
	}
	_, err := s.DB.Exec("UPDATE team_members SET role = $1 WHERE team_id = $2 AND user_id = $3", newRole, teamID, userID)
	return err
}

// DeleteTeam удаляет команду и всех её участников.
func (s *TeamStore) DeleteTeam(teamID int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec("DELETE FROM team_members WHERE team_id = $1", teamID); err != nil {
		return err
	}
	if _, err = tx.Exec("DELETE FROM teams WHERE id = $1", teamID); err != nil {
		return err
	}

	return tx.Commit()
}

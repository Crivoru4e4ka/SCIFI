package services

import (
	"errors"
	"project-MVP/db"
	"project-MVP/models"
	"strings"
)

func GetUserTeams(userID int) ([]models.Team, error) {
	query := `
		SELECT t.id, t.name, t.description, t.created_by, 
		(SELECT COUNT(*) FROM team_members tm WHERE tm.team_id = t.id) as members_count
		FROM teams t
		JOIN team_members tm ON t.id = tm.team_id
		WHERE tm.user_id = $1`

	rows, err := db.DB.Query(query, userID)
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

func CreateTeam(t models.Team) (models.Team, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return t, err
	}
	defer tx.Rollback()

	// 1. Создаем команду
	err = tx.QueryRow(`INSERT INTO teams (name, description, created_by) VALUES ($1, $2, $3) RETURNING id`,
		t.Name, t.Description, t.CreatedBy).Scan(&t.ID)
	if err != nil {
		return t, err
	}

	// 2. Добавляем создателя с проектной ролью project_lead
	_, err = tx.Exec(`INSERT INTO team_members (team_id, user_id, role) VALUES ($1, $2, 'project_lead')`, t.ID, t.CreatedBy)
	if err != nil {
		return t, err
	}

	// 3. Добавляем приглашенных по почте с ролью researcher по умолчанию
	for _, email := range t.MemberEmails {
		var invitedID int
		// Ищем ID пользователя
		err := tx.QueryRow("SELECT id FROM users WHERE email = $1", strings.TrimSpace(email)).Scan(&invitedID)
		if err != nil {
			continue // Если почта не найдена, просто идем дальше
		}
		if invitedID == t.CreatedBy {
			continue
		}

		// Записываем в связку
		_, err = tx.Exec(`INSERT INTO team_members (team_id, user_id, role) VALUES ($1, $2, 'researcher')`, t.ID, invitedID)
		if err != nil {
			return t, err
		}
	}

	err = tx.Commit()
	return t, nil
}

func GetTeamMembers(teamID int) ([]models.TeamMemberInfo, error) {
	query := `
        SELECT u.id, u.full_name, u.email, tm.role 
        FROM users u
        JOIN team_members tm ON u.id = tm.user_id
        WHERE tm.team_id = $1`

	rows, err := db.DB.Query(query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []models.TeamMemberInfo{} // инициализируем пустым массивом
	for rows.Next() {
		var m models.TeamMemberInfo
		if err := rows.Scan(&m.UserID, &m.FullName, &m.Email, &m.Role); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

// Удалить участника из команды
func RemoveMemberFromTeam(teamID int, userID int) error {
	_, err := db.DB.Exec("DELETE FROM team_members WHERE team_id = $1 AND user_id = $2", teamID, userID)
	return err
}

// Обновить информацию о команде
func UpdateTeam(teamID int, name, description string) error {
	_, err := db.DB.Exec("UPDATE teams SET name = $1, description = $2 WHERE id = $3", name, description, teamID)
	return err
}

// Добавить одного участника в существующую команду по Email
func AddMemberToTeam(teamID int, email string) (models.TeamMemberInfo, error) {
	var m models.TeamMemberInfo
	var userID int

	// 1. Ищем пользователя
	err := db.DB.QueryRow("SELECT id, full_name, email FROM users WHERE email = $1", strings.TrimSpace(email)).
		Scan(&userID, &m.FullName, &m.Email)
	if err != nil {
		return m, err
	}

	// 2. Добавляем в связку с проектной ролью researcher по умолчанию
	_, err = db.DB.Exec("INSERT INTO team_members (team_id, user_id, role) VALUES ($1, $2, 'researcher') ON CONFLICT DO NOTHING", teamID, userID)

	m.UserID = userID
	m.Role = "researcher"
	return m, err
}

// Изменить роль участника команды (допускаются только проектные роли)
func UpdateMemberRole(teamID, userID int, newRole string) error {
	newRole = strings.ToLower(strings.TrimSpace(newRole))
	if !IsValidProjectRole(newRole) {
		return errors.New("invalid team member role: must be a valid project role")
	}
	_, err := db.DB.Exec("UPDATE team_members SET role = $1 WHERE team_id = $2 AND user_id = $3", newRole, teamID, userID)
	return err
}

// Удалить команду полностью (вместе с участниками)
func DeleteTeam(teamID int) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Сначала удаляем участников
	if _, err = tx.Exec("DELETE FROM team_members WHERE team_id = $1", teamID); err != nil {
		return err
	}
	// Потом команду
	if _, err = tx.Exec("DELETE FROM teams WHERE id = $1", teamID); err != nil {
		return err
	}

	return tx.Commit()
}

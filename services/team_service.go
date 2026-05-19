package services

import (
	"project-MVP/db"
	"project-MVP/models"
)

// DefaultTeamStore — глобальный инстанс TeamStore для обратной совместимости.
var DefaultTeamStore = NewTeamStore(db.DB)

// GetTeamByID обёртка над DefaultTeamStore.
func GetTeamByID(teamID int) (models.Team, error) {
	return DefaultTeamStore.GetTeamByID(teamID)
}

// GetUserTeams обёртка над DefaultTeamStore.
func GetUserTeams(userID int) ([]models.Team, error) {
	return DefaultTeamStore.GetUserTeams(userID)
}

// CreateTeam обёртка над DefaultTeamStore.
func CreateTeam(t models.Team) (models.Team, error) {
	return DefaultTeamStore.CreateTeam(t)
}

// GetTeamMembers обёртка над DefaultTeamStore.
func GetTeamMembers(teamID int) ([]models.TeamMemberInfo, error) {
	return DefaultTeamStore.GetTeamMembers(teamID)
}

// RemoveMemberFromTeam обёртка над DefaultTeamStore.
func RemoveMemberFromTeam(teamID int, userID int) error {
	return DefaultTeamStore.RemoveMemberFromTeam(teamID, userID)
}

// UpdateTeam обёртка над DefaultTeamStore.
func UpdateTeam(teamID int, name, description string) error {
	return DefaultTeamStore.UpdateTeam(teamID, name, description)
}

// AddMemberToTeam обёртка над DefaultTeamStore.
func AddMemberToTeam(teamID int, email string) (models.TeamMemberInfo, error) {
	return DefaultTeamStore.AddMemberToTeam(teamID, email)
}

// UpdateMemberRole обёртка над DefaultTeamStore.
func UpdateMemberRole(teamID, userID int, newRole string) error {
	return DefaultTeamStore.UpdateMemberRole(teamID, userID, newRole)
}

// DeleteTeam обёртка над DefaultTeamStore.
func DeleteTeam(teamID int) error {
	return DefaultTeamStore.DeleteTeam(teamID)
}

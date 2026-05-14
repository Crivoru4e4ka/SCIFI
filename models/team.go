package models

type Team struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	CreatedBy    int      `json:"created_by"`
	MemberEmails []string `json:"member_emails"` // Список почт участников
	MembersCount int      `json:"members_count"`
}

type TeamMemberInfo struct {
	UserID   int    `json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

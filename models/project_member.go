package models

type ProjectMember struct {
	Id        int    `json:"id"`
	ProjectId int    `json:"project_id"`
	UserId    int    `json:"user_id"`
	Role      string `json:"role"`
	RoleId    *int   `json:"role_id,omitempty"`
}

package models

import "time"

// Role представляет роль (системную или проектную)
type Role struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsSystem    bool      `json:"is_system"`
	CreatedAt   time.Time `json:"created_at"`
}

// Permission представляет право доступа
type Permission struct {
	Id       int    `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

// RolePermission связь роли и права
type RolePermission struct {
	RoleId       int `json:"role_id"`
	PermissionId int `json:"permission_id"`
}

// ProjectMemberWithRole расширенная информация об участнике с ролью
type ProjectMemberWithRole struct {
	Id       int    `json:"id"`
	UserId   int    `json:"user_id"`
	ProjectId int   `json:"project_id"`
	RoleId   *int   `json:"role_id,omitempty"`
	RoleName string `json:"role_name,omitempty"`
	RoleDesc string `json:"role_description,omitempty"`
	UserName string `json:"user_name,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
}

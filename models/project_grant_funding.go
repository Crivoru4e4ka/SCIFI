package models

import "time"

// ProjectGrantFunding связывает гранты и проекты (many-to-many)
type ProjectGrantFunding struct {
	Id               int       `json:"id"`
	GrantId          int       `json:"grant_id"`
	ProjectId        int       `json:"project_id"`
	SectionId        *int      `json:"section_id,omitempty"`
	AllocatedAmount  float64   `json:"allocated_amount"`
	FundingPurpose   string    `json:"funding_purpose"`
	FundingStartDate *string   `json:"funding_start_date,omitempty"`
	FundingEndDate   *string   `json:"funding_end_date,omitempty"`
	Notes            string    `json:"notes"`
	CreatedAt        time.Time `json:"created_at"`
}

// ProjectGrantInfo расширенная информация о финансировании проекта
type ProjectGrantInfo struct {
	ProjectGrantFunding
	ProjectName string `json:"project_name"`
	ProjectKey  string `json:"project_key"`
	GrantTitle  string `json:"grant_title"`
	GrantCode   string `json:"grant_code"`
}

// GrantProjectInfo расширенная информация о проектах, финансируемых грантом
type GrantProjectInfo struct {
	ProjectGrantFunding
	ProjectName string `json:"project_name"`
	ProjectKey  string `json:"project_key"`
}

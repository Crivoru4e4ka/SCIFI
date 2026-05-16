package models

import "time"

// GrantType типы грантов
const (
	GrantTypeState         = "state"         // государственный
	GrantTypeUniversity    = "university"    // университетский
	GrantTypeInternational = "international" // международный
	GrantTypeCorporate     = "corporate"     // корпоративный
	GrantTypeInternal      = "internal"      // внутренний
)

// GrantStatus статусы грантов
const (
	GrantStatusDraft       = "draft"        // черновик
	GrantStatusSubmitted   = "submitted"    // подан
	GrantStatusUnderReview = "under_review" // на рассмотрении
	GrantStatusApproved    = "approved"     // одобрен
	GrantStatusRejected    = "rejected"     // отклонен
	GrantStatusActive      = "active"       // активен
	GrantStatusCompleted   = "completed"    // завершен
	GrantStatusSuspended   = "suspended"    // приостановлен
)

// Grant представляет научный грант
type Grant struct {
	Id                      int       `json:"id"`
	Title                   string    `json:"title"`
	Code                    string    `json:"code"`
	FundingOrganization     string    `json:"funding_organization"`
	Country                 string    `json:"country"`
	Description             string    `json:"description"`
	ScientificDirection     string    `json:"scientific_direction"`
	GrantType               string    `json:"grant_type"`
	Status                  string    `json:"status"`
	TotalAmount             float64   `json:"total_amount"`
	Currency                string    `json:"currency"`
	StartDate               *string   `json:"start_date,omitempty"`
	EndDate                 *string   `json:"end_date,omitempty"`
	ApplicationDeadline     *string   `json:"application_deadline,omitempty"`
	PrincipalInvestigatorId int       `json:"principal_investigator_id"`
	CreatedBy               int       `json:"created_by"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               *string   `json:"updated_at,omitempty"`
}

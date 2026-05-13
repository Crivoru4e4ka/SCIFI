package models

type Sprint struct {
	ID        int     `json:"id"`
	ProjectID int     `json:"project_id"`
	Name      string  `json:"name"`
	Status    string  `json:"status"`
	StartDate *string `json:"start_date"` // используем указатель, так как может быть null
	EndDate   *string `json:"end_date"`
	Goal      string  `json:"goal"`
}

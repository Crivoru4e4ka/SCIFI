package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"project-MVP/db"
	"project-MVP/models"
)

var allowedGrantTypes = map[string]bool{
	"state":         true,
	"university":    true,
	"international": true,
	"corporate":     true,
	"internal":      true,
}

var allowedGrantStatuses = map[string]bool{
	"draft":        true,
	"submitted":    true,
	"under_review": true,
	"approved":     true,
	"rejected":     true,
	"active":       true,
	"completed":    true,
	"suspended":    true,
}

func grantTypeIsValid(t string) bool {
	return allowedGrantTypes[strings.ToLower(strings.TrimSpace(t))]
}

func grantStatusIsValid(s string) bool {
	return allowedGrantStatuses[strings.ToLower(strings.TrimSpace(s))]
}

// parseDatePtr парсит строку даты формата YYYY-MM-DD в *time.Time
func parseDatePtr(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// nullTimeStr конвертирует sql.NullTime в *string (YYYY-MM-DD)
func nullTimeStr(nt sql.NullTime) *string {
	if nt.Valid {
		s := nt.Time.Format("2006-01-02")
		return &s
	}
	return nil
}

// CreateGrant создает новый грант
func CreateGrant(g models.Grant) (models.Grant, error) {
	g.Title = strings.TrimSpace(g.Title)
	g.FundingOrganization = strings.TrimSpace(g.FundingOrganization)
	g.Code = strings.TrimSpace(g.Code)
	g.ScientificDirection = strings.TrimSpace(g.ScientificDirection)
	g.Description = strings.TrimSpace(g.Description)
	g.Currency = strings.TrimSpace(g.Currency)
	if g.Currency == "" {
		g.Currency = "RUB"
	}
	g.GrantType = strings.ToLower(strings.TrimSpace(g.GrantType))
	if g.GrantType == "" {
		g.GrantType = "state"
	}
	if !grantTypeIsValid(g.GrantType) {
		return models.Grant{}, errors.New("invalid grant type")
	}
	g.Status = strings.ToLower(strings.TrimSpace(g.Status))
	if g.Status == "" {
		g.Status = "draft"
	}
	if !grantStatusIsValid(g.Status) {
		return models.Grant{}, errors.New("invalid grant status")
	}
	if g.Title == "" {
		return models.Grant{}, errors.New("grant title is required")
	}
	if g.FundingOrganization == "" {
		return models.Grant{}, errors.New("funding organization is required")
	}
	if g.CreatedBy <= 0 {
		return models.Grant{}, errors.New("created_by is required")
	}
	if g.TotalAmount < 0 {
		return models.Grant{}, errors.New("total amount cannot be negative")
	}

	startDate, err := parseDatePtr(g.StartDate)
	if err != nil {
		return models.Grant{}, errors.New("invalid start date format")
	}
	endDate, err := parseDatePtr(g.EndDate)
	if err != nil {
		return models.Grant{}, errors.New("invalid end date format")
	}
	appDeadline, err := parseDatePtr(g.ApplicationDeadline)
	if err != nil {
		return models.Grant{}, errors.New("invalid application deadline format")
	}

	if startDate != nil && endDate != nil && endDate.Before(*startDate) {
		return models.Grant{}, errors.New("end date cannot be before start date")
	}
	if appDeadline != nil && startDate != nil && appDeadline.After(*startDate) {
		return models.Grant{}, errors.New("application deadline cannot be after start date")
	}

	if _, err := GetUserByID(g.CreatedBy); err != nil {
		return models.Grant{}, err
	}
	if g.PrincipalInvestigatorId > 0 {
		if _, err := GetUserByID(g.PrincipalInvestigatorId); err != nil {
			return models.Grant{}, err
		}
	}

	now := time.Now()
	query := `INSERT INTO grants (
		title, code, funding_organization, country, description,
		scientific_direction, grant_type, status, total_amount, currency,
		start_date, end_date, application_deadline,
		principal_investigator_id, created_by, created_at, updated_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
	RETURNING id`

	var newID int
	err = db.DB.QueryRow(query,
		g.Title, g.Code, g.FundingOrganization, g.Country, g.Description,
		g.ScientificDirection, g.GrantType, g.Status, g.TotalAmount, g.Currency,
		startDate, endDate, appDeadline,
		nullInt(g.PrincipalInvestigatorId), g.CreatedBy, now, now,
	).Scan(&newID)
	if err != nil {
		return models.Grant{}, err
	}

	g.Id = newID
	g.CreatedAt = now
	return g, nil
}

// GetGrantByID возвращает грант по ID
func GetGrantByID(id int) (models.Grant, error) {
	var g models.Grant
	var code, country, desc, sciDir, currency sql.NullString
	var startDate, endDate, appDeadline, updatedAt sql.NullTime
	var piID, createdBy sql.NullInt64

	query := `SELECT id, title, code, funding_organization, country, description,
		scientific_direction, grant_type, status, total_amount, currency,
		start_date, end_date, application_deadline,
		principal_investigator_id, created_by, created_at, updated_at
	FROM grants WHERE id = $1`

	err := db.DB.QueryRow(query, id).Scan(
		&g.Id, &g.Title, &code, &g.FundingOrganization, &country, &desc,
		&sciDir, &g.GrantType, &g.Status, &g.TotalAmount, &currency,
		&startDate, &endDate, &appDeadline,
		&piID, &createdBy, &g.CreatedAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return g, ErrNotFound
		}
		return g, err
	}

	g.Code = nullStringValue(code)
	g.Country = nullStringValue(country)
	g.Description = nullStringValue(desc)
	g.ScientificDirection = nullStringValue(sciDir)
	g.Currency = nullStringValue(currency)
	g.StartDate = nullTimeStr(startDate)
	g.EndDate = nullTimeStr(endDate)
	g.ApplicationDeadline = nullTimeStr(appDeadline)
	g.PrincipalInvestigatorId = int(nullInt64Value(piID))
	g.CreatedBy = int(nullInt64Value(createdBy))
	g.UpdatedAt = nullTimeStr(updatedAt)
	return g, nil
}

// UpdateGrant обновляет грант
func UpdateGrant(id int, g models.Grant) (models.Grant, error) {
	existing, err := GetGrantByID(id)
	if err != nil {
		return models.Grant{}, err
	}

	g.Title = strings.TrimSpace(g.Title)
	g.FundingOrganization = strings.TrimSpace(g.FundingOrganization)
	g.Code = strings.TrimSpace(g.Code)
	g.ScientificDirection = strings.TrimSpace(g.ScientificDirection)
	g.Description = strings.TrimSpace(g.Description)
	g.Currency = strings.TrimSpace(g.Currency)

	if g.Title == "" {
		return models.Grant{}, errors.New("grant title is required")
	}
	if g.FundingOrganization == "" {
		return models.Grant{}, errors.New("funding organization is required")
	}
	if g.GrantType != "" && !grantTypeIsValid(g.GrantType) {
		return models.Grant{}, errors.New("invalid grant type")
	}
	if g.Status != "" && !grantStatusIsValid(g.Status) {
		return models.Grant{}, errors.New("invalid grant status")
	}
	if g.TotalAmount < 0 {
		return models.Grant{}, errors.New("total amount cannot be negative")
	}

	startDate, err := parseDatePtr(g.StartDate)
	if err != nil {
		return models.Grant{}, errors.New("invalid start date format")
	}
	endDate, err := parseDatePtr(g.EndDate)
	if err != nil {
		return models.Grant{}, errors.New("invalid end date format")
	}
	appDeadline, err := parseDatePtr(g.ApplicationDeadline)
	if err != nil {
		return models.Grant{}, errors.New("invalid application deadline format")
	}

	if startDate != nil && endDate != nil && endDate.Before(*startDate) {
		return models.Grant{}, errors.New("end date cannot be before start date")
	}

	if g.GrantType == "" {
		g.GrantType = existing.GrantType
	}
	if g.Status == "" {
		g.Status = existing.Status
	}
	if g.Currency == "" {
		g.Currency = existing.Currency
	}

	query := `UPDATE grants SET
		title = $1, code = $2, funding_organization = $3, country = $4, description = $5,
		scientific_direction = $6, grant_type = $7, status = $8, total_amount = $9, currency = $10,
		start_date = $11, end_date = $12, application_deadline = $13,
		principal_investigator_id = $14, updated_at = $15
	WHERE id = $16`

	now := time.Now()
	_, err = db.DB.Exec(query,
		g.Title, g.Code, g.FundingOrganization, g.Country, g.Description,
		g.ScientificDirection, g.GrantType, g.Status, g.TotalAmount, g.Currency,
		startDate, endDate, appDeadline,
		nullInt(g.PrincipalInvestigatorId), now, id,
	)
	if err != nil {
		return models.Grant{}, err
	}

	return GetGrantByID(id)
}

// DeleteGrant удаляет грант и связанные распределения (CASCADE)
func DeleteGrant(id int) error {
	res, err := db.DB.Exec(`DELETE FROM grants WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// ListGrantsFilters параметры фильтрации грантов
type ListGrantsFilters struct {
	Status       string
	GrantType    string
	Organization string
	Search       string
	DateFrom     *time.Time
	DateTo       *time.Time
	Limit        int
	Offset       int
}

// ListGrants возвращает список грантов с фильтрацией
func ListGrants(f ListGrantsFilters) ([]models.Grant, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}
	argIdx := 1

	if f.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, strings.ToLower(f.Status))
		argIdx++
	}
	if f.GrantType != "" {
		where = append(where, fmt.Sprintf("grant_type = $%d", argIdx))
		args = append(args, strings.ToLower(f.GrantType))
		argIdx++
	}
	if f.Organization != "" {
		where = append(where, fmt.Sprintf("funding_organization ILIKE $%d", argIdx))
		args = append(args, "%"+f.Organization+"%")
		argIdx++
	}
	if f.Search != "" {
		where = append(where, fmt.Sprintf("(title ILIKE $%d OR code ILIKE $%d OR description ILIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, "%"+f.Search+"%")
		argIdx++
	}
	if f.DateFrom != nil {
		where = append(where, fmt.Sprintf("start_date >= $%d", argIdx))
		args = append(args, *f.DateFrom)
		argIdx++
	}
	if f.DateTo != nil {
		where = append(where, fmt.Sprintf("end_date <= $%d", argIdx))
		args = append(args, *f.DateTo)
		argIdx++
	}

	whereClause := strings.Join(where, " AND ")

	// Считаем общее количество
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM grants WHERE %s", whereClause)
	if err := db.DB.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Основной запрос
	query := fmt.Sprintf(`SELECT id, title, code, funding_organization, country, description,
		scientific_direction, grant_type, status, total_amount, currency,
		start_date, end_date, application_deadline,
		principal_investigator_id, created_by, created_at, updated_at
	FROM grants WHERE %s ORDER BY created_at DESC`, whereClause)

	if f.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, f.Limit)
		argIdx++
	}
	if f.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, f.Offset)
		argIdx++
	}

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var grants []models.Grant
	for rows.Next() {
		var g models.Grant
		var code, country, desc, sciDir, currency sql.NullString
		var startDate, endDate, appDeadline, updatedAt sql.NullTime
		var piID, createdBy sql.NullInt64

		err := rows.Scan(
			&g.Id, &g.Title, &code, &g.FundingOrganization, &country, &desc,
			&sciDir, &g.GrantType, &g.Status, &g.TotalAmount, &currency,
			&startDate, &endDate, &appDeadline,
			&piID, &createdBy, &g.CreatedAt, &updatedAt,
		)
		if err != nil {
			log.Printf("Scan error in ListGrants: %v", err)
			continue
		}

		g.Code = nullStringValue(code)
		g.Country = nullStringValue(country)
		g.Description = nullStringValue(desc)
		g.ScientificDirection = nullStringValue(sciDir)
		g.Currency = nullStringValue(currency)
		g.StartDate = nullTimeStr(startDate)
		g.EndDate = nullTimeStr(endDate)
		g.ApplicationDeadline = nullTimeStr(appDeadline)
		g.PrincipalInvestigatorId = int(nullInt64Value(piID))
		g.CreatedBy = int(nullInt64Value(createdBy))
		g.UpdatedAt = nullTimeStr(updatedAt)
		grants = append(grants, g)
	}
	return grants, total, nil
}

// GetGrantTotalAllocated возвращает сумму всех распределений по гранту
func GetGrantTotalAllocated(grantID int) (float64, error) {
	var total sql.NullFloat64
	err := db.DB.QueryRow(`SELECT COALESCE(SUM(allocated_amount), 0) FROM project_grant_funding WHERE grant_id = $1`, grantID).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total.Float64, nil
}

// GetGrantProjects возвращает проекты, финансируемые грантом
func GetGrantProjects(grantID int) ([]models.GrantProjectInfo, error) {
	query := `
		SELECT
			pgf.id, pgf.grant_id, pgf.project_id, pgf.section_id,
			pgf.allocated_amount, pgf.funding_purpose, pgf.funding_start_date,
			pgf.funding_end_date, pgf.notes, pgf.created_at,
			p.name, p.key
		FROM project_grant_funding pgf
		JOIN projects p ON p.id = pgf.project_id
		WHERE pgf.grant_id = $1
		ORDER BY pgf.created_at DESC`

	rows, err := db.DB.Query(query, grantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.GrantProjectInfo
	for rows.Next() {
		var info models.GrantProjectInfo
		var sectionID sql.NullInt64
		var fStart, fEnd sql.NullTime
		var purpose, notes sql.NullString

		err := rows.Scan(
			&info.Id, &info.GrantId, &info.ProjectId, &sectionID,
			&info.AllocatedAmount, &purpose, &fStart,
			&fEnd, &notes, &info.CreatedAt,
			&info.ProjectName, &info.ProjectKey,
		)
		if err != nil {
			log.Printf("Scan error in GetGrantProjects: %v", err)
			continue
		}
		if sectionID.Valid {
			v := int(sectionID.Int64)
			info.SectionId = &v
		}
		info.FundingPurpose = nullStringValue(purpose)
		info.FundingStartDate = nullTimeStr(fStart)
		info.FundingEndDate = nullTimeStr(fEnd)
		info.Notes = nullStringValue(notes)
		list = append(list, info)
	}
	return list, nil
}

// Вспомогательные функции
func nullStringValue(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func nullInt64Value(ni sql.NullInt64) int64 {
	if ni.Valid {
		return ni.Int64
	}
	return 0
}

func nullInt(v int) interface{} {
	if v > 0 {
		return v
	}
	return nil
}

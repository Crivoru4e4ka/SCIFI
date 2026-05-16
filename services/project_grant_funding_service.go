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

// CreateProjectGrantFunding создает связь финансирования между грантом и проектом
func CreateProjectGrantFunding(f models.ProjectGrantFunding) (models.ProjectGrantFunding, error) {
	f.FundingPurpose = strings.TrimSpace(f.FundingPurpose)
	f.Notes = strings.TrimSpace(f.Notes)

	if f.GrantId <= 0 {
		return models.ProjectGrantFunding{}, errors.New("grant_id is required")
	}
	if f.ProjectId <= 0 {
		return models.ProjectGrantFunding{}, errors.New("project_id is required")
	}
	if f.AllocatedAmount < 0 {
		return models.ProjectGrantFunding{}, errors.New("allocated amount cannot be negative")
	}

	fundingStart, err := parseDatePtr(f.FundingStartDate)
	if err != nil {
		return models.ProjectGrantFunding{}, errors.New("invalid funding start date format")
	}
	fundingEnd, err := parseDatePtr(f.FundingEndDate)
	if err != nil {
		return models.ProjectGrantFunding{}, errors.New("invalid funding end date format")
	}
	if fundingStart != nil && fundingEnd != nil && fundingEnd.Before(*fundingStart) {
		return models.ProjectGrantFunding{}, errors.New("funding end date cannot be before start date")
	}

	// Проверяем существование гранта
	grant, err := GetGrantByID(f.GrantId)
	if err != nil {
		return models.ProjectGrantFunding{}, err
	}

	// Проверяем существование проекта
	if _, err := GetProjectByID(f.ProjectId); err != nil {
		return models.ProjectGrantFunding{}, err
	}

	// Проверяем, что сумма распределений не превышает бюджет
	if f.AllocatedAmount > 0 {
		totalAllocated, err := GetGrantTotalAllocated(f.GrantId)
		if err != nil {
			return models.ProjectGrantFunding{}, err
		}
		if totalAllocated+f.AllocatedAmount > grant.TotalAmount {
			return models.ProjectGrantFunding{}, fmt.Errorf("allocated amount (%.2f) exceeds grant budget (%.2f). Already allocated: %.2f", f.AllocatedAmount, grant.TotalAmount, totalAllocated)
		}
	}

	query := `INSERT INTO project_grant_funding (
		grant_id, project_id, section_id, allocated_amount,
		funding_purpose, funding_start_date, funding_end_date, notes, created_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	RETURNING id`

	now := time.Now()
	var newID int
	err = db.DB.QueryRow(query,
		f.GrantId, f.ProjectId, f.SectionId, f.AllocatedAmount,
		f.FundingPurpose, fundingStart, fundingEnd, f.Notes, now,
	).Scan(&newID)
	if err != nil {
		if strings.Contains(err.Error(), "unique_grant_project_section") {
			return models.ProjectGrantFunding{}, errors.New("funding for this project and section already exists")
		}
		return models.ProjectGrantFunding{}, err
	}

	f.Id = newID
	f.CreatedAt = now
	return f, nil
}

// UpdateProjectGrantFunding обновляет распределение
func UpdateProjectGrantFunding(id int, f models.ProjectGrantFunding) (models.ProjectGrantFunding, error) {
	existing, err := GetProjectGrantFundingByID(id)
	if err != nil {
		return models.ProjectGrantFunding{}, err
	}

	f.FundingPurpose = strings.TrimSpace(f.FundingPurpose)
	f.Notes = strings.TrimSpace(f.Notes)

	if f.AllocatedAmount < 0 {
		return models.ProjectGrantFunding{}, errors.New("allocated amount cannot be negative")
	}

	fundingStart, err := parseDatePtr(f.FundingStartDate)
	if err != nil {
		return models.ProjectGrantFunding{}, errors.New("invalid funding start date format")
	}
	fundingEnd, err := parseDatePtr(f.FundingEndDate)
	if err != nil {
		return models.ProjectGrantFunding{}, errors.New("invalid funding end date format")
	}

	// Пересчитываем бюджет
	if f.AllocatedAmount != existing.AllocatedAmount {
		grant, err := GetGrantByID(existing.GrantId)
		if err != nil {
			return models.ProjectGrantFunding{}, err
		}
		totalAllocated, err := GetGrantTotalAllocated(existing.GrantId)
		if err != nil {
			return models.ProjectGrantFunding{}, err
		}
		newTotal := totalAllocated - existing.AllocatedAmount + f.AllocatedAmount
		if newTotal > grant.TotalAmount {
			return models.ProjectGrantFunding{}, fmt.Errorf("updated allocated amount would exceed grant budget (%.2f). Current total: %.2f", grant.TotalAmount, totalAllocated)
		}
	}

	query := `UPDATE project_grant_funding SET
		section_id = $1, allocated_amount = $2, funding_purpose = $3,
		funding_start_date = $4, funding_end_date = $5, notes = $6
	WHERE id = $7`

	_, err = db.DB.Exec(query,
		f.SectionId, f.AllocatedAmount, f.FundingPurpose,
		fundingStart, fundingEnd, f.Notes, id,
	)
	if err != nil {
		return models.ProjectGrantFunding{}, err
	}
	return GetProjectGrantFundingByID(id)
}

// DeleteProjectGrantFunding удаляет распределение
func DeleteProjectGrantFunding(id int) error {
	res, err := db.DB.Exec(`DELETE FROM project_grant_funding WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// GetProjectGrantFundingByID возвращает распределение по ID
func GetProjectGrantFundingByID(id int) (models.ProjectGrantFunding, error) {
	var f models.ProjectGrantFunding
	var sectionID sql.NullInt64
	var fStart, fEnd sql.NullTime
	var purpose, notes sql.NullString

	query := `SELECT id, grant_id, project_id, section_id, allocated_amount,
		funding_purpose, funding_start_date, funding_end_date, notes, created_at
	FROM project_grant_funding WHERE id = $1`

	err := db.DB.QueryRow(query, id).Scan(
		&f.Id, &f.GrantId, &f.ProjectId, &sectionID, &f.AllocatedAmount,
		&purpose, &fStart, &fEnd, &notes, &f.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return f, ErrNotFound
		}
		return f, err
	}
	if sectionID.Valid {
		v := int(sectionID.Int64)
		f.SectionId = &v
	}
	f.FundingPurpose = nullStringValue(purpose)
	f.FundingStartDate = nullTimeStr(fStart)
	f.FundingEndDate = nullTimeStr(fEnd)
	f.Notes = nullStringValue(notes)
	return f, nil
}

// GetProjectGrants возвращает гранты, финансирующие проект
func GetProjectGrants(projectID int) ([]models.ProjectGrantInfo, error) {
	query := `
		SELECT
			pgf.id, pgf.grant_id, pgf.project_id, pgf.section_id,
			pgf.allocated_amount, pgf.funding_purpose, pgf.funding_start_date,
			pgf.funding_end_date, pgf.notes, pgf.created_at,
			p.name, p.key,
			g.title, g.code
		FROM project_grant_funding pgf
		JOIN projects p ON p.id = pgf.project_id
		JOIN grants g ON g.id = pgf.grant_id
		WHERE pgf.project_id = $1
		ORDER BY pgf.created_at DESC`

	rows, err := db.DB.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.ProjectGrantInfo
	for rows.Next() {
		var info models.ProjectGrantInfo
		var sectionID sql.NullInt64
		var fStart, fEnd sql.NullTime
		var purpose, notes sql.NullString

		err := rows.Scan(
			&info.Id, &info.GrantId, &info.ProjectId, &sectionID,
			&info.AllocatedAmount, &purpose, &fStart,
			&fEnd, &notes, &info.CreatedAt,
			&info.ProjectName, &info.ProjectKey,
			&info.GrantTitle, &info.GrantCode,
		)
		if err != nil {
			log.Printf("Scan error in GetProjectGrants: %v", err)
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

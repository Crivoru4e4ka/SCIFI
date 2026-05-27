package services

import (
	"database/sql"

	"project-MVP/db"
	"project-MVP/models"
)

// DatasetStore инкапсулирует операции с датасетами.
type DatasetStore struct {
	DB db.DBPool
}

// NewDatasetStore создает новый экземпляр DatasetStore.
func NewDatasetStore(database db.DBPool) *DatasetStore {
	return &DatasetStore{DB: database}
}

// GetAllDatasets возвращает все датасеты.
func (s *DatasetStore) GetAllDatasets() ([]models.Dataset, error) {
	rows, err := s.DB.Query(`SELECT id, project_id, name, description, version, data_url, parameters, created_by, created_at FROM datasets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Dataset
	for rows.Next() {
		var ds models.Dataset
		var parameters []byte
		var createdBy sql.NullInt64
		if err := rows.Scan(&ds.ID, &ds.ProjectID, &ds.Name, &ds.Description, &ds.Version, &ds.DataURL, &parameters, &createdBy, &ds.CreatedAt); err != nil {
			return nil, err
		}
		ds.Parameters = parameters
		if createdBy.Valid {
			ds.CreatedBy = int(createdBy.Int64)
		}
		result = append(result, ds)
	}
	return result, nil
}

// GetDatasetsByProject возвращает датасеты проекта.
func (s *DatasetStore) GetDatasetsByProject(projectID int) ([]models.Dataset, error) {
	rows, err := s.DB.Query(`SELECT id, project_id, name, description, version, data_url, parameters, created_by, created_at FROM datasets WHERE project_id = $1 ORDER BY created_at DESC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Dataset
	for rows.Next() {
		var ds models.Dataset
		var parameters []byte
		var createdBy sql.NullInt64
		if err := rows.Scan(&ds.ID, &ds.ProjectID, &ds.Name, &ds.Description, &ds.Version, &ds.DataURL, &parameters, &createdBy, &ds.CreatedAt); err != nil {
			return nil, err
		}
		ds.Parameters = parameters
		if createdBy.Valid {
			ds.CreatedBy = int(createdBy.Int64)
		}
		result = append(result, ds)
	}
	return result, nil
}

// GetDatasetByID возвращает датасет по ID.
func (s *DatasetStore) GetDatasetByID(id int) (models.Dataset, error) {
	var ds models.Dataset
	var parameters []byte
	var createdBy sql.NullInt64
	err := s.DB.QueryRow(`SELECT id, project_id, name, description, version, data_url, parameters, created_by, created_at FROM datasets WHERE id = $1`, id).
		Scan(&ds.ID, &ds.ProjectID, &ds.Name, &ds.Description, &ds.Version, &ds.DataURL, &parameters, &createdBy, &ds.CreatedAt)
	if err != nil {
		return ds, err
	}
	ds.Parameters = parameters
	if createdBy.Valid {
		ds.CreatedBy = int(createdBy.Int64)
	}
	return ds, nil
}

// CreateDataset создает новый датасет.
func (s *DatasetStore) CreateDataset(ds models.Dataset) (models.Dataset, error) {
	row := s.DB.QueryRow(`INSERT INTO datasets (project_id, name, description, version, data_url, parameters, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`,
		ds.ProjectID, ds.Name, ds.Description, ds.Version, ds.DataURL, ds.Parameters, ds.CreatedBy)
	var id int
	var createdAt sql.NullTime
	if err := row.Scan(&id, &createdAt); err != nil {
		return ds, err
	}
	ds.ID = id
	if createdAt.Valid {
		ds.CreatedAt = createdAt.Time
	}
	return ds, nil
}

// DeleteDataset удаляет датасет.
func (s *DatasetStore) DeleteDataset(id int) error {
	_, err := s.DB.Exec(`DELETE FROM datasets WHERE id = $1`, id)
	return err
}

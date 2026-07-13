package repositories

import (
	"database/sql"
	"fmt"

	"github.com/isaac/statescore/internal/models"
)

// DataSourceRepository handles data source database operations.
type DataSourceRepository struct {
	db *sql.DB
}

// NewDataSourceRepository creates a new DataSourceRepository.
func NewDataSourceRepository(db *sql.DB) *DataSourceRepository {
	return &DataSourceRepository{db: db}
}

// List returns all data sources.
func (r *DataSourceRepository) List() ([]models.DataSource, error) {
	rows, err := r.db.Query(
		`SELECT id, name, publisher, source_url, license, format, description, created_at, updated_at FROM data_sources ORDER BY name`,
	)
	if err != nil {
		return nil, fmt.Errorf("list data sources: %w", err)
	}
	defer rows.Close()

	var sources []models.DataSource
	for rows.Next() {
		var s models.DataSource
		if err := rows.Scan(&s.ID, &s.Name, &s.Publisher, &s.SourceURL, &s.License, &s.Format, &s.Description, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan data source: %w", err)
		}
		sources = append(sources, s)
	}
	return sources, rows.Err()
}

// GetByID returns a data source by its ID.
func (r *DataSourceRepository) GetByID(id int64) (*models.DataSource, error) {
	var s models.DataSource
	err := r.db.QueryRow(
		`SELECT id, name, publisher, source_url, license, format, description, created_at, updated_at FROM data_sources WHERE id = ?`,
		id,
	).Scan(&s.ID, &s.Name, &s.Publisher, &s.SourceURL, &s.License, &s.Format, &s.Description, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get data source: %w", err)
	}
	return &s, nil
}

// Create inserts a new data source.
func (r *DataSourceRepository) Create(s *models.DataSource) error {
	result, err := r.db.Exec(
		`INSERT INTO data_sources (name, publisher, source_url, license, format, description) VALUES (?, ?, ?, ?, ?, ?)`,
		s.Name, s.Publisher, s.SourceURL, s.License, s.Format, s.Description,
	)
	if err != nil {
		return fmt.Errorf("create data source: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get data source id: %w", err)
	}
	s.ID = id
	return nil
}

// Update updates an existing data source.
func (r *DataSourceRepository) Update(s *models.DataSource) error {
	_, err := r.db.Exec(
		`UPDATE data_sources SET name = ?, publisher = ?, source_url = ?, license = ?, format = ?, description = ?, updated_at = datetime('now') WHERE id = ?`,
		s.Name, s.Publisher, s.SourceURL, s.License, s.Format, s.Description, s.ID,
	)
	if err != nil {
		return fmt.Errorf("update data source: %w", err)
	}
	return nil
}

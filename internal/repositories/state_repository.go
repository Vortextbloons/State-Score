package repositories

import (
	"database/sql"
	"fmt"

	"github.com/isaac/statescore/internal/models"
)

// StateRepository handles state database operations.
type StateRepository struct {
	db *sql.DB
}

// NewStateRepository creates a new StateRepository.
func NewStateRepository(db *sql.DB) *StateRepository {
	return &StateRepository{db: db}
}

// List returns all states, optionally filtered by region.
func (r *StateRepository) List(region string) ([]models.State, error) {
	query := `SELECT id, code, name, region, division, population, population_year, population_source_id, created_at, updated_at FROM states`
	args := []any{}

	if region != "" {
		query += ` WHERE region = ?`
		args = append(args, region)
	}

	query += ` ORDER BY name`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list states: %w", err)
	}
	defer rows.Close()

	var states []models.State
	for rows.Next() {
		var s models.State
		if err := rows.Scan(&s.ID, &s.Code, &s.Name, &s.Region, &s.Division, &s.Population, &s.PopulationYear, &s.PopulationSourceID, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan state: %w", err)
		}
		states = append(states, s)
	}
	return states, rows.Err()
}

// GetByCode returns a state by its two-letter code.
func (r *StateRepository) GetByCode(code string) (*models.State, error) {
	var s models.State
	err := r.db.QueryRow(
		`SELECT id, code, name, region, division, population, population_year, population_source_id, created_at, updated_at FROM states WHERE code = ?`,
		code,
	).Scan(&s.ID, &s.Code, &s.Name, &s.Region, &s.Division, &s.Population, &s.PopulationYear, &s.PopulationSourceID, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get state by code: %w", err)
	}
	return &s, nil
}

// GetByID returns a state by its ID.
func (r *StateRepository) GetByID(id int64) (*models.State, error) {
	var s models.State
	err := r.db.QueryRow(
		`SELECT id, code, name, region, division, population, population_year, population_source_id, created_at, updated_at FROM states WHERE id = ?`,
		id,
	).Scan(&s.ID, &s.Code, &s.Name, &s.Region, &s.Division, &s.Population, &s.PopulationYear, &s.PopulationSourceID, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get state by id: %w", err)
	}
	return &s, nil
}

// Regions returns all distinct regions.
func (r *StateRepository) Regions() ([]string, error) {
	rows, err := r.db.Query(`SELECT DISTINCT region FROM states WHERE region IS NOT NULL ORDER BY region`)
	if err != nil {
		return nil, fmt.Errorf("list regions: %w", err)
	}
	defer rows.Close()

	var regions []string
	for rows.Next() {
		var region string
		if err := rows.Scan(&region); err != nil {
			return nil, fmt.Errorf("scan region: %w", err)
		}
		regions = append(regions, region)
	}
	return regions, rows.Err()
}

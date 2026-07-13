package repositories

import (
	"database/sql"
	"fmt"

	"github.com/isaac/statescore/internal/models"
)

// MetricRepository handles metric database operations.
type MetricRepository struct {
	db *sql.DB
}

// NewMetricRepository creates a new MetricRepository.
func NewMetricRepository(db *sql.DB) *MetricRepository {
	return &MetricRepository{db: db}
}

// List returns all active metrics, optionally filtered by category.
func (r *MetricRepository) List(categoryID int64) ([]models.Metric, error) {
	query := `SELECT id, category_id, slug, name, description, unit, higher_is_better, normalization_method, default_weight, source_id, active, created_at, updated_at FROM metrics WHERE active = 1`
	args := []any{}

	if categoryID > 0 {
		query += ` AND category_id = ?`
		args = append(args, categoryID)
	}

	query += ` ORDER BY category_id, name`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list metrics: %w", err)
	}
	defer rows.Close()

	var metrics []models.Metric
	for rows.Next() {
		var m models.Metric
		var sourceID sql.NullInt64
		if err := rows.Scan(
			&m.ID, &m.CategoryID, &m.Slug, &m.Name, &m.Description, &m.Unit,
			&m.HigherIsBetter, &m.NormalizationMethod, &m.DefaultWeight,
			&sourceID, &m.Active, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan metric: %w", err)
		}
		if sourceID.Valid {
			m.SourceID = &sourceID.Int64
		}
		metrics = append(metrics, m)
	}
	return metrics, rows.Err()
}

// GetBySlug returns a metric by its slug.
func (r *MetricRepository) GetBySlug(slug string) (*models.Metric, error) {
	var m models.Metric
	var sourceID sql.NullInt64
	err := r.db.QueryRow(
		`SELECT id, category_id, slug, name, description, unit, higher_is_better, normalization_method, default_weight, source_id, active, created_at, updated_at FROM metrics WHERE slug = ?`,
		slug,
	).Scan(
		&m.ID, &m.CategoryID, &m.Slug, &m.Name, &m.Description, &m.Unit,
		&m.HigherIsBetter, &m.NormalizationMethod, &m.DefaultWeight,
		&sourceID, &m.Active, &m.CreatedAt, &m.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get metric by slug: %w", err)
	}
	if sourceID.Valid {
		m.SourceID = &sourceID.Int64
	}
	return &m, nil
}

// GetByID returns a metric by its ID.
func (r *MetricRepository) GetByID(id int64) (*models.Metric, error) {
	var m models.Metric
	var sourceID sql.NullInt64
	err := r.db.QueryRow(
		`SELECT id, category_id, slug, name, description, unit, higher_is_better, normalization_method, default_weight, source_id, active, created_at, updated_at FROM metrics WHERE id = ?`,
		id,
	).Scan(
		&m.ID, &m.CategoryID, &m.Slug, &m.Name, &m.Description, &m.Unit,
		&m.HigherIsBetter, &m.NormalizationMethod, &m.DefaultWeight,
		&sourceID, &m.Active, &m.CreatedAt, &m.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get metric by id: %w", err)
	}
	if sourceID.Valid {
		m.SourceID = &sourceID.Int64
	}
	return &m, nil
}

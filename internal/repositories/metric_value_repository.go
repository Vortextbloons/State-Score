package repositories

import (
	"database/sql"
	"fmt"

	"github.com/isaac/statescore/internal/models"
)

// MetricValueRepository handles metric value database operations.
type MetricValueRepository struct {
	db *sql.DB
}

// NewMetricValueRepository creates a new MetricValueRepository.
func NewMetricValueRepository(db *sql.DB) *MetricValueRepository {
	return &MetricValueRepository{db: db}
}

// ListByState returns all metric values for a state, optionally filtered by year.
func (r *MetricValueRepository) ListByState(stateID int64, year int) ([]models.MetricValue, error) {
	query := `SELECT id, state_id, metric_id, year, value, source_record_id, import_id, created_at FROM metric_values WHERE state_id = ?`
	args := []any{stateID}

	if year > 0 {
		query += ` AND year = ?`
		args = append(args, year)
	}

	query += ` ORDER BY metric_id, year`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list metric values: %w", err)
	}
	defer rows.Close()

	values := make([]models.MetricValue, 0)
	for rows.Next() {
		var mv models.MetricValue
		var importID sql.NullInt64
		if err := rows.Scan(&mv.ID, &mv.StateID, &mv.MetricID, &mv.Year, &mv.Value, &mv.SourceRecordID, &importID, &mv.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan metric value: %w", err)
		}
		if importID.Valid {
			mv.ImportID = &importID.Int64
		}
		values = append(values, mv)
	}
	return values, rows.Err()
}

// ListAll returns metric values for all states, optionally filtered by year.
// It is the bulk counterpart to ListByState for rankings and exports.
func (r *MetricValueRepository) ListAll(year int) ([]models.MetricValue, error) {
	query := `SELECT id, state_id, metric_id, year, value, source_record_id, import_id, created_at FROM metric_values`
	args := []any{}
	if year > 0 {
		query += ` WHERE year = ?`
		args = append(args, year)
	}
	query += ` ORDER BY state_id, metric_id, year`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list all metric values: %w", err)
	}
	defer rows.Close()

	values := make([]models.MetricValue, 0)
	for rows.Next() {
		var mv models.MetricValue
		var importID sql.NullInt64
		if err := rows.Scan(&mv.ID, &mv.StateID, &mv.MetricID, &mv.Year, &mv.Value, &mv.SourceRecordID, &importID, &mv.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan metric value: %w", err)
		}
		if importID.Valid {
			mv.ImportID = &importID.Int64
		}
		values = append(values, mv)
	}
	return values, rows.Err()
}

// ListByMetric returns all values for a metric, optionally filtered by year.
func (r *MetricValueRepository) ListByMetric(metricID int64, year int) ([]models.MetricValue, error) {
	query := `SELECT id, state_id, metric_id, year, value, source_record_id, import_id, created_at FROM metric_values WHERE metric_id = ?`
	args := []any{metricID}

	if year > 0 {
		query += ` AND year = ?`
		args = append(args, year)
	}

	query += ` ORDER BY state_id, year`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list metric values: %w", err)
	}
	defer rows.Close()

	values := make([]models.MetricValue, 0)
	for rows.Next() {
		var mv models.MetricValue
		var importID sql.NullInt64
		if err := rows.Scan(&mv.ID, &mv.StateID, &mv.MetricID, &mv.Year, &mv.Value, &mv.SourceRecordID, &importID, &mv.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan metric value: %w", err)
		}
		if importID.Valid {
			mv.ImportID = &importID.Int64
		}
		values = append(values, mv)
	}
	return values, rows.Err()
}

// Get returns a specific metric value for a state, metric, and year.
func (r *MetricValueRepository) Get(stateID, metricID int64, year int) (*models.MetricValue, error) {
	var mv models.MetricValue
	var importID sql.NullInt64
	err := r.db.QueryRow(
		`SELECT id, state_id, metric_id, year, value, source_record_id, import_id, created_at FROM metric_values WHERE state_id = ? AND metric_id = ? AND year = ?`,
		stateID, metricID, year,
	).Scan(&mv.ID, &mv.StateID, &mv.MetricID, &mv.Year, &mv.Value, &mv.SourceRecordID, &importID, &mv.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get metric value: %w", err)
	}
	if importID.Valid {
		mv.ImportID = &importID.Int64
	}
	return &mv, nil
}

// Upsert inserts or updates a metric value.
func (r *MetricValueRepository) Upsert(mv *models.MetricValue) error {
	_, err := r.db.Exec(
		`INSERT INTO metric_values (state_id, metric_id, year, value, source_record_id, import_id)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON CONFLICT (state_id, metric_id, year, import_id) DO UPDATE SET value = ?, source_record_id = ?`,
		mv.StateID, mv.MetricID, mv.Year, mv.Value, mv.SourceRecordID, mv.ImportID,
		mv.Value, mv.SourceRecordID,
	)
	if err != nil {
		return fmt.Errorf("upsert metric value: %w", err)
	}
	return nil
}

// AvailableYears returns all years that have metric values.
func (r *MetricValueRepository) AvailableYears() ([]int, error) {
	rows, err := r.db.Query(`SELECT DISTINCT year FROM metric_values ORDER BY year DESC`)
	if err != nil {
		return nil, fmt.Errorf("list years: %w", err)
	}
	defer rows.Close()

	var years []int
	for rows.Next() {
		var year int
		if err := rows.Scan(&year); err != nil {
			return nil, fmt.Errorf("scan year: %w", err)
		}
		years = append(years, year)
	}
	return years, rows.Err()
}

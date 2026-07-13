package repositories

import (
	"database/sql"
	"fmt"

	"github.com/isaac/statescore/internal/models"
)

// ImportRepository handles import database operations.
type ImportRepository struct {
	db *sql.DB
}

// NewImportRepository creates a new ImportRepository.
func NewImportRepository(db *sql.DB) *ImportRepository {
	return &ImportRepository{db: db}
}

// List returns all imports, optionally limited.
func (r *ImportRepository) List(limit int) ([]models.Import, error) {
	query := `SELECT id, source_id, status, started_at, completed_at, records_read, records_inserted, records_rejected, checksum, error_summary FROM imports ORDER BY started_at DESC`
	if limit > 0 {
		query += fmt.Sprintf(` LIMIT %d`, limit)
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("list imports: %w", err)
	}
	defer rows.Close()

	var imports []models.Import
	for rows.Next() {
		var i models.Import
		var sourceID sql.NullInt64
		if err := rows.Scan(&i.ID, &sourceID, &i.Status, &i.StartedAt, &i.CompletedAt, &i.RecordsRead, &i.RecordsInserted, &i.RecordsRejected, &i.Checksum, &i.ErrorSummary); err != nil {
			return nil, fmt.Errorf("scan import: %w", err)
		}
		if sourceID.Valid {
			i.SourceID = &sourceID.Int64
		}
		imports = append(imports, i)
	}
	return imports, rows.Err()
}

// GetByID returns an import by its ID.
func (r *ImportRepository) GetByID(id int64) (*models.Import, error) {
	var i models.Import
	var sourceID sql.NullInt64
	err := r.db.QueryRow(
		`SELECT id, source_id, status, started_at, completed_at, records_read, records_inserted, records_rejected, checksum, error_summary FROM imports WHERE id = ?`,
		id,
	).Scan(&i.ID, &sourceID, &i.Status, &i.StartedAt, &i.CompletedAt, &i.RecordsRead, &i.RecordsInserted, &i.RecordsRejected, &i.Checksum, &i.ErrorSummary)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get import: %w", err)
	}
	if sourceID.Valid {
		i.SourceID = &sourceID.Int64
	}
	return &i, nil
}

// Create inserts a new import.
func (r *ImportRepository) Create(i *models.Import) error {
	result, err := r.db.Exec(
		`INSERT INTO imports (source_id, status) VALUES (?, ?)`,
		i.SourceID, i.Status,
	)
	if err != nil {
		return fmt.Errorf("create import: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get import id: %w", err)
	}
	i.ID = id
	return nil
}

// Update updates an existing import.
func (r *ImportRepository) Update(i *models.Import) error {
	_, err := r.db.Exec(
		`UPDATE imports SET status = ?, started_at = ?, completed_at = ?, records_read = ?, records_inserted = ?, records_rejected = ?, checksum = ?, error_summary = ? WHERE id = ?`,
		i.Status, i.StartedAt, i.CompletedAt, i.RecordsRead, i.RecordsInserted, i.RecordsRejected, i.Checksum, i.ErrorSummary, i.ID,
	)
	if err != nil {
		return fmt.Errorf("update import: %w", err)
	}
	return nil
}

// ListErrors returns all errors for an import.
func (r *ImportRepository) ListErrors(importID int64) ([]models.ImportError, error) {
	rows, err := r.db.Query(
		`SELECT id, import_id, row_number, field_name, raw_value, error_message FROM import_errors WHERE import_id = ? ORDER BY row_number`,
		importID,
	)
	if err != nil {
		return nil, fmt.Errorf("list import errors: %w", err)
	}
	defer rows.Close()

	var errors []models.ImportError
	for rows.Next() {
		var e models.ImportError
		var rowNumber sql.NullInt64
		if err := rows.Scan(&e.ID, &e.ImportID, &rowNumber, &e.FieldName, &e.RawValue, &e.ErrorMessage); err != nil {
			return nil, fmt.Errorf("scan import error: %w", err)
		}
		if rowNumber.Valid {
			rowNum := int(rowNumber.Int64)
			e.RowNumber = &rowNum
		}
		errors = append(errors, e)
	}
	return errors, rows.Err()
}

// AddError adds an error to an import.
func (r *ImportRepository) AddError(e *models.ImportError) error {
	_, err := r.db.Exec(
		`INSERT INTO import_errors (import_id, row_number, field_name, raw_value, error_message) VALUES (?, ?, ?, ?, ?)`,
		e.ImportID, e.RowNumber, e.FieldName, e.RawValue, e.ErrorMessage,
	)
	if err != nil {
		return fmt.Errorf("add import error: %w", err)
	}
	return nil
}

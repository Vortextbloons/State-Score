package repositories

import (
	"database/sql"
	"fmt"

	"github.com/isaac/statescore/internal/models"
)

// ProfileRepository handles scoring profile database operations.
type ProfileRepository struct {
	db *sql.DB
}

// NewProfileRepository creates a new ProfileRepository.
func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// List returns all scoring profiles.
func (r *ProfileRepository) List() ([]models.ScoringProfile, error) {
	rows, err := r.db.Query(
		`SELECT id, name, description, is_default, is_system, created_at, updated_at FROM scoring_profiles ORDER BY name`,
	)
	if err != nil {
		return nil, fmt.Errorf("list profiles: %w", err)
	}
	defer rows.Close()

	var profiles []models.ScoringProfile
	for rows.Next() {
		var p models.ScoringProfile
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.IsDefault, &p.IsSystem, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan profile: %w", err)
		}
		profiles = append(profiles, p)
	}
	return profiles, rows.Err()
}

// GetByID returns a profile by its ID.
func (r *ProfileRepository) GetByID(id int64) (*models.ScoringProfile, error) {
	var p models.ScoringProfile
	err := r.db.QueryRow(
		`SELECT id, name, description, is_default, is_system, created_at, updated_at FROM scoring_profiles WHERE id = ?`,
		id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.IsDefault, &p.IsSystem, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}
	return &p, nil
}

// GetDefault returns the default scoring profile.
func (r *ProfileRepository) GetDefault() (*models.ScoringProfile, error) {
	var p models.ScoringProfile
	err := r.db.QueryRow(
		`SELECT id, name, description, is_default, is_system, created_at, updated_at FROM scoring_profiles WHERE is_default = 1`,
	).Scan(&p.ID, &p.Name, &p.Description, &p.IsDefault, &p.IsSystem, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get default profile: %w", err)
	}
	return &p, nil
}

// Create inserts a new scoring profile.
func (r *ProfileRepository) Create(p *models.ScoringProfile) error {
	result, err := r.db.Exec(
		`INSERT INTO scoring_profiles (name, description, is_default, is_system) VALUES (?, ?, ?, ?)`,
		p.Name, p.Description, p.IsDefault, p.IsSystem,
	)
	if err != nil {
		return fmt.Errorf("create profile: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get profile id: %w", err)
	}
	p.ID = id
	return nil
}

// Update updates an existing scoring profile.
func (r *ProfileRepository) Update(p *models.ScoringProfile) error {
	_, err := r.db.Exec(
		`UPDATE scoring_profiles SET name = ?, description = ?, is_default = ?, updated_at = datetime('now') WHERE id = ?`,
		p.Name, p.Description, p.IsDefault, p.ID,
	)
	if err != nil {
		return fmt.Errorf("update profile: %w", err)
	}
	return nil
}

// Delete deletes a scoring profile.
func (r *ProfileRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM scoring_profiles WHERE id = ? AND is_system = 0`, id)
	if err != nil {
		return fmt.Errorf("delete profile: %w", err)
	}
	return nil
}

// GetCategoryWeights returns category weights for a profile.
func (r *ProfileRepository) GetCategoryWeights(profileID int64) ([]models.ProfileCategoryWeight, error) {
	rows, err := r.db.Query(
		`SELECT profile_id, category_id, weight FROM profile_category_weights WHERE profile_id = ?`,
		profileID,
	)
	if err != nil {
		return nil, fmt.Errorf("list category weights: %w", err)
	}
	defer rows.Close()

	var weights []models.ProfileCategoryWeight
	for rows.Next() {
		var w models.ProfileCategoryWeight
		if err := rows.Scan(&w.ProfileID, &w.CategoryID, &w.Weight); err != nil {
			return nil, fmt.Errorf("scan category weight: %w", err)
		}
		weights = append(weights, w)
	}
	return weights, rows.Err()
}

// SetCategoryWeights sets category weights for a profile.
func (r *ProfileRepository) SetCategoryWeights(profileID int64, weights []models.ProfileCategoryWeight) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete existing weights
	if _, err := tx.Exec(`DELETE FROM profile_category_weights WHERE profile_id = ?`, profileID); err != nil {
		return fmt.Errorf("delete category weights: %w", err)
	}

	// Insert new weights
	for _, w := range weights {
		if _, err := tx.Exec(
			`INSERT INTO profile_category_weights (profile_id, category_id, weight) VALUES (?, ?, ?)`,
			profileID, w.CategoryID, w.Weight,
		); err != nil {
			return fmt.Errorf("insert category weight: %w", err)
		}
	}

	return tx.Commit()
}

// GetMetricWeights returns metric weights for a profile.
func (r *ProfileRepository) GetMetricWeights(profileID int64) ([]models.ProfileMetricWeight, error) {
	rows, err := r.db.Query(
		`SELECT profile_id, metric_id, weight FROM profile_metric_weights WHERE profile_id = ?`,
		profileID,
	)
	if err != nil {
		return nil, fmt.Errorf("list metric weights: %w", err)
	}
	defer rows.Close()

	var weights []models.ProfileMetricWeight
	for rows.Next() {
		var w models.ProfileMetricWeight
		if err := rows.Scan(&w.ProfileID, &w.MetricID, &w.Weight); err != nil {
			return nil, fmt.Errorf("scan metric weight: %w", err)
		}
		weights = append(weights, w)
	}
	return weights, rows.Err()
}

// SetMetricWeights sets metric weights for a profile.
func (r *ProfileRepository) SetMetricWeights(profileID int64, weights []models.ProfileMetricWeight) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete existing weights
	if _, err := tx.Exec(`DELETE FROM profile_metric_weights WHERE profile_id = ?`, profileID); err != nil {
		return fmt.Errorf("delete metric weights: %w", err)
	}

	// Insert new weights
	for _, w := range weights {
		if _, err := tx.Exec(
			`INSERT INTO profile_metric_weights (profile_id, metric_id, weight) VALUES (?, ?, ?)`,
			profileID, w.MetricID, w.Weight,
		); err != nil {
			return fmt.Errorf("insert metric weight: %w", err)
		}
	}

	return tx.Commit()
}

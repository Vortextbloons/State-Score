package repositories

import (
	"database/sql"
	"fmt"

	"github.com/isaac/statescore/internal/models"
)

// CategoryRepository handles category database operations.
type CategoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository creates a new CategoryRepository.
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// List returns all categories ordered by display_order.
func (r *CategoryRepository) List() ([]models.Category, error) {
	rows, err := r.db.Query(
		`SELECT id, slug, name, description, default_weight, display_order FROM categories ORDER BY display_order`,
	)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Slug, &c.Name, &c.Description, &c.DefaultWeight, &c.DisplayOrder); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

// GetBySlug returns a category by its slug.
func (r *CategoryRepository) GetBySlug(slug string) (*models.Category, error) {
	var c models.Category
	err := r.db.QueryRow(
		`SELECT id, slug, name, description, default_weight, display_order FROM categories WHERE slug = ?`,
		slug,
	).Scan(&c.ID, &c.Slug, &c.Name, &c.Description, &c.DefaultWeight, &c.DisplayOrder)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get category by slug: %w", err)
	}
	return &c, nil
}

// GetByID returns a category by its ID.
func (r *CategoryRepository) GetByID(id int64) (*models.Category, error) {
	var c models.Category
	err := r.db.QueryRow(
		`SELECT id, slug, name, description, default_weight, display_order FROM categories WHERE id = ?`,
		id,
	).Scan(&c.ID, &c.Slug, &c.Name, &c.Description, &c.DefaultWeight, &c.DisplayOrder)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get category by id: %w", err)
	}
	return &c, nil
}

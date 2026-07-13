package database_test

import (
	"path/filepath"
	"testing"

	"github.com/isaac/statescore/internal/database"
)

func TestOpenAndMigrate(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "statescore.db")
	db, err := database.Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if err := database.Migrate(db); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	// Second run must be a no-op.
	if err := database.Migrate(db); err != nil {
		t.Fatalf("Migrate second pass: %v", err)
	}

	var version string
	if err := db.QueryRow(`SELECT version FROM schema_migrations`).Scan(&version); err != nil {
		t.Fatalf("query schema_migrations: %v", err)
	}
	if version != "000001_initial_schema" {
		t.Fatalf("version = %q, want 000001_initial_schema", version)
	}

	tables := []string{
		"states",
		"categories",
		"metrics",
		"metric_values",
		"data_sources",
		"imports",
		"import_errors",
		"scoring_profiles",
		"profile_category_weights",
		"profile_metric_weights",
		"score_snapshots",
		"category_score_snapshots",
		"application_settings",
	}
	for _, table := range tables {
		var name string
		err := db.QueryRow(
			`SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?`,
			table,
		).Scan(&name)
		if err != nil {
			t.Fatalf("missing table %s: %v", table, err)
		}
	}
}

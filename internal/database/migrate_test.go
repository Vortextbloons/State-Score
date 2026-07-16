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

	var applied int
	if err := db.QueryRow(`SELECT count(*) FROM schema_migrations`).Scan(&applied); err != nil {
		t.Fatalf("query schema_migrations: %v", err)
	}
	if applied < 10 {
		t.Fatalf("applied migrations = %d, want at least 10", applied)
	}
	var latest string
	if err := db.QueryRow(`SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1`).Scan(&latest); err != nil {
		t.Fatalf("latest migration: %v", err)
	}
	if latest != "000010_add_priority_metrics" {
		t.Fatalf("latest version = %q, want 000010_add_priority_metrics", latest)
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
		"metric_value_quality",
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

func TestBundledMetricData(t *testing.T) {
	t.Parallel()

	db, err := database.Open(filepath.Join(t.TempDir(), "statescore.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := database.Migrate(db); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	for _, metric := range []struct {
		slug string
		year int
		want int
	}{
		{slug: "life-expectancy", year: 2022, want: 50},
		{slug: "cost-of-living-index", year: 2024, want: 50},
		{slug: "unemployment-rate", year: 2024, want: 50},
		{slug: "median-household-income", year: 2024, want: 50},
		{slug: "high-school-graduation-rate", year: 2024, want: 50},
		{slug: "bachelors-degree-attainment", year: 2024, want: 50},
		{slug: "annual-employment-growth", year: 2024, want: 50},
		{slug: "adult-obesity-prevalence", year: 2024, want: 49},
		{slug: "property-crime-rate", year: 2024, want: 50},
		{slug: "young-adult-college-enrollment", year: 2024, want: 50},
		{slug: "renter-housing-cost-burden", year: 2024, want: 50},
	} {
		var count int
		err := db.QueryRow(`SELECT count(*) FROM metric_values mv JOIN metrics m ON m.id=mv.metric_id WHERE m.slug=? AND mv.year=?`, metric.slug, metric.year).Scan(&count)
		if err != nil {
			t.Fatalf("count %s: %v", metric.slug, err)
		}
		if count != metric.want {
			t.Fatalf("%s rows = %d, want %d", metric.slug, count, metric.want)
		}
	}

	for _, slug := range []string{
		"uninsured-rate",
		"median-rent",
	} {
		var active int
		err := db.QueryRow(`SELECT active FROM metrics WHERE slug=?`, slug).Scan(&active)
		if err != nil {
			t.Fatalf("active %s: %v", slug, err)
		}
		if active != 0 {
			t.Fatalf("%s active = %d, want 0 until data is seeded", slug, active)
		}
	}

	var excluded int
	if err := db.QueryRow(`SELECT count(*) FROM metric_value_quality WHERE scoring_eligible=0`).Scan(&excluded); err != nil || excluded != 6 {
		t.Fatalf("excluded crime observations = %d, err=%v; want 6", excluded, err)
	}

	var name string
	if err := db.QueryRow(`SELECT name FROM metrics WHERE slug='cost-of-living-index'`).Scan(&name); err != nil {
		t.Fatalf("rpp name: %v", err)
	}
	if name != "Regional price parity" {
		t.Fatalf("cost-of-living-index name = %q, want Regional price parity", name)
	}
}

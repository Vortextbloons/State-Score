package scoring_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/isaac/statescore/internal/database"
	"github.com/isaac/statescore/internal/repositories"
	"github.com/isaac/statescore/internal/scoring"
)

func TestRecalculateUsesAsOfYearAndActiveMetrics(t *testing.T) {
	t.Parallel()

	db, err := database.Open(filepath.Join(t.TempDir(), "statescore.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := database.Migrate(db); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	var active int
	if err := db.QueryRow(`SELECT count(*) FROM metrics WHERE active=1`).Scan(&active); err != nil {
		t.Fatalf("count active metrics: %v", err)
	}
	if active != 18 {
		t.Fatalf("active metrics = %d, want 18", active)
	}

	profile, err := repositories.NewProfileRepository(db).GetDefault()
	if err != nil || profile == nil {
		t.Fatalf("default profile: %v", err)
	}

	count, err := scoring.Recalculate(context.Background(), db, profile.ID, 2024)
	if err != nil {
		t.Fatalf("Recalculate: %v", err)
	}
	if count != 50 {
		t.Fatalf("states calculated = %d, want 50", count)
	}

	rows, err := repositories.NewScoreRepository(db).ListByProfileYear(profile.ID, 2024, scoring.CalculationVersion)
	if err != nil {
		t.Fatalf("ListByProfileYear: %v", err)
	}
	if len(rows) != 50 {
		t.Fatalf("snapshots = %d, want 50", len(rows))
	}
	incomplete := 0
	for _, row := range rows {
		if row.Snapshot.Completeness < 1 {
			incomplete++
		}
		if len(row.Categories) != 5 {
			t.Fatalf("state %d categories = %d, want 5", row.Snapshot.StateID, len(row.Categories))
		}
	}
	if incomplete != 3 {
		t.Fatalf("incomplete states = %d, want 3 states without a coverage-qualified FBI fallback", incomplete)
	}
}

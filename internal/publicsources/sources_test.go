package publicsources

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/isaac/statescore/internal/database"
)

type fakeAdapter struct{}

func (fakeAdapter) Spec() Spec {
	return Spec{ID: "fake", Name: "Fake source", Publisher: "Test", MetricSlugs: []string{"young-adult-college-enrollment"}, DefaultYear: 2024, Available: true}
}
func (fakeAdapter) SourceName() string { return "ACS 2024 Subject Table S1401" }
func (fakeAdapter) Fetch(context.Context, int) (Batch, error) {
	coverage := 100.0
	return Batch{[]Observation{{StateCode: "AL", MetricSlug: "young-adult-college-enrollment", SourceRecordID: "fixture", Year: 2024, Value: 42, Quality: &Quality{ReportingCoverage: &coverage, ScoringEligible: true}}}}, nil
}

func TestAdapterRegistryAndAtomicPersistence(t *testing.T) {
	db, err := database.Open(filepath.Join(t.TempDir(), "refresh.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err = database.Migrate(db); err != nil {
		t.Fatal(err)
	}
	service := NewService(db, NewRegistry(fakeAdapter{}))
	if got := service.Specs(); len(got) != 1 || got[0].ID != "fake" {
		t.Fatalf("specs=%v", got)
	}
	importID, err := service.Prepare("fake")
	if err != nil {
		t.Fatal(err)
	}
	if err = service.Run(context.Background(), "fake", 2024, importID); err != nil {
		t.Fatal(err)
	}
	var status string
	var count int
	if err = db.QueryRow(`SELECT status FROM imports WHERE id=?`, importID).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != "completed" {
		t.Fatalf("status=%s", status)
	}
	if err = db.QueryRow(`SELECT count(*) FROM metric_values WHERE import_id=?`, importID).Scan(&count); err != nil || count != 1 {
		t.Fatalf("count=%d err=%v", count, err)
	}
}

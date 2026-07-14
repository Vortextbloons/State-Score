package scoring

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/isaac/statescore/internal/models"
	"github.com/isaac/statescore/internal/repositories"
)

// Recalculate creates canonical snapshots for every state for a profile / as-of year.
// For each active metric it uses the latest completed import value with year ≤ asOfYear,
// so staggered release calendars still produce one coherent composite.
func Recalculate(ctx context.Context, db *sql.DB, profileID int64, asOfYear int) (int, error) {
	states, err := repositories.NewStateRepository(db).List("")
	if err != nil {
		return 0, err
	}
	categories, err := repositories.NewCategoryRepository(db).List()
	if err != nil {
		return 0, err
	}
	metrics, err := repositories.NewMetricRepository(db).List(0)
	if err != nil {
		return 0, err
	}
	profile := repositories.NewProfileRepository(db)
	cw, err := profile.GetCategoryWeights(profileID)
	if err != nil {
		return 0, err
	}
	mw, err := profile.GetMetricWeights(profileID)
	if err != nil {
		return 0, err
	}
	categoryWeight := map[int64]float64{}
	for _, w := range cw {
		categoryWeight[w.CategoryID] = w.Weight
	}
	metricWeight := map[int64]float64{}
	for _, w := range mw {
		metricWeight[w.MetricID] = w.Weight
	}

	metricInputs := make([]MetricInput, 0, len(metrics))
	activeCategories := map[int64]struct{}{}
	for _, m := range metrics {
		obs, err := loadAsOfObservations(ctx, db, m.ID, asOfYear)
		if err != nil {
			return 0, err
		}
		scores, err := Normalize(obs, NormalizationMethod(m.NormalizationMethod), m.HigherIsBetter, nil, nil)
		if err != nil {
			return 0, fmt.Errorf("normalize %s: %w", m.Slug, err)
		}
		w := m.DefaultWeight
		if v, ok := metricWeight[m.ID]; ok {
			w = v
		}
		metricInputs = append(metricInputs, MetricInput{m.ID, m.CategoryID, w, scores})
		activeCategories[m.CategoryID] = struct{}{}
	}

	categoryInputs := make([]CategoryInput, 0, len(categories))
	for _, c := range categories {
		if _, ok := activeCategories[c.ID]; !ok {
			continue
		}
		w := c.DefaultWeight
		if v, ok := categoryWeight[c.ID]; ok {
			w = v
		}
		categoryInputs = append(categoryInputs, CategoryInput{c.ID, w})
	}
	if len(categoryInputs) == 0 {
		return 0, nil
	}

	ids := make([]int64, len(states))
	for i, s := range states {
		ids[i] = s.ID
	}
	calculated, err := Calculate(ids, categoryInputs, metricInputs)
	if err != nil {
		return 0, err
	}
	repo := repositories.NewScoreRepository(db)
	for _, score := range calculated {
		cats := make([]models.CategoryScoreSnapshot, 0, len(score.Categories))
		for id, c := range score.Categories {
			cats = append(cats, models.CategoryScoreSnapshot{CategoryID: id, Score: c.Score, Completeness: c.Completeness})
		}
		snap := models.ScoreSnapshot{
			ProfileID:          profileID,
			StateID:            score.StateID,
			Year:               asOfYear,
			OverallScore:       score.Overall.Score,
			Completeness:       score.Overall.Completeness,
			CalculatedAt:       time.Now().UTC(),
			CalculationVersion: CalculationVersion,
		}
		if err := repo.Save(&snap, cats); err != nil {
			return 0, err
		}
	}
	return len(calculated), nil
}

func loadAsOfObservations(ctx context.Context, db *sql.DB, metricID int64, asOfYear int) ([]Observation, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT mv.state_id, mv.year, mv.id, mv.value
		FROM metric_values mv
		JOIN imports i ON i.id = mv.import_id
		WHERE mv.metric_id = ?
		  AND mv.year <= ?
		  AND i.status IN ('completed', 'completed_with_errors')
		ORDER BY mv.state_id ASC, mv.year DESC, mv.id DESC`, metricID, asOfYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var obs []Observation
	seen := map[int64]struct{}{}
	for rows.Next() {
		var stateID, rowID int64
		var year int
		var value float64
		if err = rows.Scan(&stateID, &year, &rowID, &value); err != nil {
			return nil, err
		}
		if _, ok := seen[stateID]; ok {
			continue
		}
		seen[stateID] = struct{}{}
		v := value
		obs = append(obs, Observation{stateID, &v})
	}
	return obs, rows.Err()
}

// EnsureLatest recalculates the default profile for the latest available as-of year.
func EnsureLatest(ctx context.Context, db *sql.DB) error {
	profile, err := repositories.NewProfileRepository(db).GetDefault()
	if err != nil {
		return err
	}
	if profile == nil {
		return nil
	}
	years, err := repositories.NewMetricValueRepository(db).AvailableYears()
	if err != nil {
		return err
	}
	if len(years) == 0 {
		return nil
	}
	_, err = Recalculate(ctx, db, profile.ID, years[0])
	return err
}

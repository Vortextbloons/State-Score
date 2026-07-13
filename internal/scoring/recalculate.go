package scoring

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/isaac/statescore/internal/models"
	"github.com/isaac/statescore/internal/repositories"
	"time"
)

// Recalculate creates canonical snapshots for every state for a profile/year.
func Recalculate(ctx context.Context, db *sql.DB, profileID int64, year int) (int, error) {
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
	categoryInputs := make([]CategoryInput, 0, len(categories))
	for _, c := range categories {
		w := c.DefaultWeight
		if v, ok := categoryWeight[c.ID]; ok {
			w = v
		}
		categoryInputs = append(categoryInputs, CategoryInput{c.ID, w})
	}
	metricInputs := make([]MetricInput, 0, len(metrics))
	for _, m := range metrics {
		rows, err := db.QueryContext(ctx, `SELECT mv.state_id,mv.value FROM metric_values mv JOIN imports i ON i.id=mv.import_id WHERE mv.metric_id=? AND mv.year=? AND i.status IN ('completed','completed_with_errors') AND mv.id=(SELECT max(mv2.id) FROM metric_values mv2 JOIN imports i2 ON i2.id=mv2.import_id WHERE mv2.state_id=mv.state_id AND mv2.metric_id=mv.metric_id AND mv2.year=mv.year AND i2.status IN ('completed','completed_with_errors'))`, m.ID, year)
		if err != nil {
			return 0, err
		}
		var obs []Observation
		for rows.Next() {
			var id int64
			var value float64
			if err = rows.Scan(&id, &value); err != nil {
				rows.Close()
				return 0, err
			}
			v := value
			obs = append(obs, Observation{id, &v})
		}
		rows.Close()
		scores, err := Normalize(obs, NormalizationMethod(m.NormalizationMethod), m.HigherIsBetter, nil, nil)
		if err != nil {
			return 0, fmt.Errorf("normalize %s: %w", m.Slug, err)
		}
		w := m.DefaultWeight
		if v, ok := metricWeight[m.ID]; ok {
			w = v
		}
		metricInputs = append(metricInputs, MetricInput{m.ID, m.CategoryID, w, scores})
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
		snap := models.ScoreSnapshot{ProfileID: profileID, StateID: score.StateID, Year: year, OverallScore: score.Overall.Score, Completeness: score.Overall.Completeness, CalculatedAt: time.Now().UTC(), CalculationVersion: CalculationVersion}
		if err := repo.Save(&snap, cats); err != nil {
			return 0, err
		}
	}
	return len(calculated), nil
}

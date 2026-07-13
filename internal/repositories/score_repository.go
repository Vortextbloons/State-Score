package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/isaac/statescore/internal/models"
)

// ScoreRepository atomically persists reproducible overall and category scores.
type ScoreRepository struct{ db *sql.DB }

func NewScoreRepository(db *sql.DB) *ScoreRepository { return &ScoreRepository{db: db} }

func (r *ScoreRepository) Save(snapshot *models.ScoreSnapshot, categories []models.CategoryScoreSnapshot) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin score transaction: %w", err)
	}
	defer tx.Rollback()
	if snapshot.CalculatedAt.IsZero() {
		snapshot.CalculatedAt = time.Now().UTC()
	}
	var id int64
	err = tx.QueryRow(`INSERT INTO score_snapshots (profile_id,state_id,year,overall_score,completeness,calculated_at,calculation_version) VALUES (?,?,?,?,?,?,?) ON CONFLICT(profile_id,state_id,year,calculation_version) DO UPDATE SET overall_score=excluded.overall_score, completeness=excluded.completeness, calculated_at=excluded.calculated_at RETURNING id`, snapshot.ProfileID, snapshot.StateID, snapshot.Year, snapshot.OverallScore, snapshot.Completeness, snapshot.CalculatedAt, snapshot.CalculationVersion).Scan(&id)
	if err != nil {
		return fmt.Errorf("save score snapshot: %w", err)
	}
	snapshot.ID = id
	if _, err = tx.Exec(`DELETE FROM category_score_snapshots WHERE score_snapshot_id=?`, id); err != nil {
		return fmt.Errorf("replace category scores: %w", err)
	}
	for _, c := range categories {
		if _, err = tx.Exec(`INSERT INTO category_score_snapshots(score_snapshot_id,category_id,score,completeness) VALUES(?,?,?,?)`, id, c.CategoryID, c.Score, c.Completeness); err != nil {
			return fmt.Errorf("save category score: %w", err)
		}
	}
	return tx.Commit()
}

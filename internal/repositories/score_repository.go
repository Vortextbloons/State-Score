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

// StateScoreRow is one state's overall snapshot plus category breakdowns.
type StateScoreRow struct {
	Snapshot   models.ScoreSnapshot
	Categories []models.CategoryScoreSnapshot
}

// ListByProfileYear returns canonical score snapshots for a profile and as-of year.
func (r *ScoreRepository) ListByProfileYear(profileID int64, year int, version string) ([]StateScoreRow, error) {
	rows, err := r.db.Query(`
		SELECT id, profile_id, state_id, year, overall_score, completeness, calculated_at, calculation_version
		FROM score_snapshots
		WHERE profile_id = ? AND year = ? AND calculation_version = ?
		ORDER BY overall_score DESC, state_id ASC`, profileID, year, version)
	if err != nil {
		return nil, fmt.Errorf("list score snapshots: %w", err)
	}

	var snaps []models.ScoreSnapshot
	for rows.Next() {
		var snap models.ScoreSnapshot
		var calculatedAt string
		if err := rows.Scan(&snap.ID, &snap.ProfileID, &snap.StateID, &snap.Year, &snap.OverallScore, &snap.Completeness, &calculatedAt, &snap.CalculationVersion); err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan score snapshot: %w", err)
		}
		if t, err := time.Parse(time.RFC3339, calculatedAt); err == nil {
			snap.CalculatedAt = t
		} else if t, err := time.Parse("2006-01-02 15:04:05", calculatedAt); err == nil {
			snap.CalculatedAt = t.UTC()
		}
		snaps = append(snaps, snap)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()

	out := make([]StateScoreRow, 0, len(snaps))
	for _, snap := range snaps {
		cats, err := r.listCategories(snap.ID)
		if err != nil {
			return nil, err
		}
		out = append(out, StateScoreRow{Snapshot: snap, Categories: cats})
	}
	return out, nil
}

func (r *ScoreRepository) listCategories(snapshotID int64) ([]models.CategoryScoreSnapshot, error) {
	rows, err := r.db.Query(`
		SELECT score_snapshot_id, category_id, score, completeness
		FROM category_score_snapshots
		WHERE score_snapshot_id = ?
		ORDER BY category_id`, snapshotID)
	if err != nil {
		return nil, fmt.Errorf("list category scores: %w", err)
	}
	defer rows.Close()

	var out []models.CategoryScoreSnapshot
	for rows.Next() {
		var c models.CategoryScoreSnapshot
		if err := rows.Scan(&c.ScoreSnapshotID, &c.CategoryID, &c.Score, &c.Completeness); err != nil {
			return nil, fmt.Errorf("scan category score: %w", err)
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// HasSnapshots reports whether any snapshots exist for the profile/year/version.
func (r *ScoreRepository) HasSnapshots(profileID int64, year int, version string) (bool, error) {
	var n int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM score_snapshots WHERE profile_id=? AND year=? AND calculation_version=?`, profileID, year, version).Scan(&n)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

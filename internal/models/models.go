package models

import "time"

// State represents a US state.
type State struct {
	ID        int64  `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	Region    string `json:"region,omitempty"`
	Division  string `json:"division,omitempty"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// Category represents a scoring category.
type Category struct {
	ID            int64   `json:"id"`
	Slug          string  `json:"slug"`
	Name          string  `json:"name"`
	Description   string  `json:"description,omitempty"`
	DefaultWeight float64 `json:"defaultWeight"`
	DisplayOrder  int     `json:"displayOrder"`
}

// Metric represents a measurable metric within a category.
type Metric struct {
	ID                  int64   `json:"id"`
	CategoryID          int64   `json:"categoryId"`
	Slug                string  `json:"slug"`
	Name                string  `json:"name"`
	Description         string  `json:"description,omitempty"`
	Unit                string  `json:"unit,omitempty"`
	HigherIsBetter      bool    `json:"higherIsBetter"`
	NormalizationMethod string  `json:"normalizationMethod"`
	DefaultWeight       float64 `json:"defaultWeight"`
	SourceID            *int64  `json:"sourceId,omitempty"`
	Active              bool    `json:"active"`
	CreatedAt           string  `json:"createdAt"`
	UpdatedAt           string  `json:"updatedAt"`
}

// MetricValue represents a single data point for a state, metric, and year.
type MetricValue struct {
	ID             int64   `json:"id"`
	StateID        int64   `json:"stateId"`
	MetricID       int64   `json:"metricId"`
	Year           int     `json:"year"`
	Value          float64 `json:"value"`
	SourceRecordID string  `json:"sourceRecordId,omitempty"`
	ImportID       *int64  `json:"importId,omitempty"`
	CreatedAt      string  `json:"createdAt"`
}

// DataSource represents an external data source.
type DataSource struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Publisher   string `json:"publisher,omitempty"`
	SourceURL   string `json:"sourceUrl,omitempty"`
	License     string `json:"license,omitempty"`
	Format      string `json:"format,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// Import represents a data import operation.
type Import struct {
	ID              int64   `json:"id"`
	SourceID        *int64  `json:"sourceId,omitempty"`
	Status          string  `json:"status"`
	StartedAt       *string `json:"startedAt,omitempty"`
	CompletedAt     *string `json:"completedAt,omitempty"`
	RecordsRead     int     `json:"recordsRead"`
	RecordsInserted int     `json:"recordsInserted"`
	RecordsRejected int     `json:"recordsRejected"`
	Checksum        string  `json:"checksum,omitempty"`
	ErrorSummary    string  `json:"errorSummary,omitempty"`
}

// ImportError represents an error that occurred during an import.
type ImportError struct {
	ID           int64  `json:"id"`
	ImportID     int64  `json:"importId"`
	RowNumber    *int   `json:"rowNumber,omitempty"`
	FieldName    string `json:"fieldName,omitempty"`
	RawValue     string `json:"rawValue,omitempty"`
	ErrorMessage string `json:"errorMessage"`
}

// ScoringProfile represents a weighting configuration.
type ScoringProfile struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsDefault   bool   `json:"isDefault"`
	IsSystem    bool   `json:"isSystem"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// ProfileCategoryWeight stores category weights for a profile.
type ProfileCategoryWeight struct {
	ProfileID  int64   `json:"profileId"`
	CategoryID int64   `json:"categoryId"`
	Weight     float64 `json:"weight"`
}

// ProfileMetricWeight stores metric weights for a profile.
type ProfileMetricWeight struct {
	ProfileID int64   `json:"profileId"`
	MetricID  int64   `json:"metricId"`
	Weight    float64 `json:"weight"`
}

// ScoreSnapshot stores a calculated overall score.
type ScoreSnapshot struct {
	ID                 int64     `json:"id"`
	ProfileID          int64     `json:"profileId"`
	StateID            int64     `json:"stateId"`
	Year               int       `json:"year"`
	OverallScore       float64   `json:"overallScore"`
	Completeness       float64   `json:"completeness"`
	CalculatedAt       time.Time `json:"calculatedAt"`
	CalculationVersion string    `json:"calculationVersion"`
}

// CategoryScoreSnapshot stores a calculated category score.
type CategoryScoreSnapshot struct {
	ScoreSnapshotID int64   `json:"scoreSnapshotId"`
	CategoryID      int64   `json:"categoryId"`
	Score           float64 `json:"score"`
	Completeness    float64 `json:"completeness"`
}

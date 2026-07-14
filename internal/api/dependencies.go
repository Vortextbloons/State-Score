package api

import (
	"database/sql"
	"github.com/isaac/statescore/internal/models"
	"github.com/isaac/statescore/internal/repositories"
)

type stateStore interface {
	List(string) ([]models.State, error)
	GetByCode(string) (*models.State, error)
	Regions() ([]string, error)
}
type categoryStore interface {
	List() ([]models.Category, error)
}
type metricStore interface {
	List(int64) ([]models.Metric, error)
	GetByID(int64) (*models.Metric, error)
}
type metricValueStore interface {
	ListByState(int64, int) ([]models.MetricValue, error)
	ListAll(int) ([]models.MetricValue, error)
	AvailableYears() ([]int, error)
}
type sourceStore interface {
	List() ([]models.DataSource, error)
	GetByID(int64) (*models.DataSource, error)
	Create(*models.DataSource) error
	Update(*models.DataSource) error
}
type importStore interface {
	List(int) ([]models.Import, error)
	GetByID(int64) (*models.Import, error)
	Create(*models.Import) error
	Update(*models.Import) error
	ListErrors(int64) ([]models.ImportError, error)
	AddError(*models.ImportError) error
}
type profileStore interface {
	List() ([]models.ScoringProfile, error)
	GetByID(int64) (*models.ScoringProfile, error)
	GetDefault() (*models.ScoringProfile, error)
	GetCategoryWeights(int64) ([]models.ProfileCategoryWeight, error)
}
type scoreStore interface {
	ListByProfileYear(int64, int, string) ([]repositories.StateScoreRow, error)
	HasSnapshots(int64, int, string) (bool, error)
}

// Dependencies is the API composition boundary for production adapters and test fakes.
type Dependencies struct {
	States       stateStore
	Categories   categoryStore
	Metrics      metricStore
	MetricValues metricValueStore
	Sources      sourceStore
	Imports      importStore
	Profiles     profileStore
	Scores       scoreStore
}

func repositoryDependencies(db *sql.DB) Dependencies {
	return Dependencies{repositories.NewStateRepository(db), repositories.NewCategoryRepository(db), repositories.NewMetricRepository(db), repositories.NewMetricValueRepository(db), repositories.NewDataSourceRepository(db), repositories.NewImportRepository(db), repositories.NewProfileRepository(db), repositories.NewScoreRepository(db)}
}

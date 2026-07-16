# State Data Model

Core data structures, schema definitions, and relationships.

## Entity Relationship

```
State ─┐
       ├──< MetricValue >── Metric >── Category
       │
DataSource ──< Import ──< ImportError
                     │
                     └──< MetricValue
ScoringProfile ──< ProfileCategoryWeight
               ──< ProfileMetricWeight
               ──< ScoreSnapshot >── CategoryScoreSnapshot
```

## Core Entities

### State

A US state or territory.

| Field | Type | Description |
|-------|------|-------------|
| `id` | `int64` | Primary key |
| `code` | `string` | Two-letter abbreviation (e.g., `CA`), UNIQUE |
| `name` | `string` | Full name |
| `region` | `string` | Census region: `South`, `West`, `Northeast`, `Midwest` |
| `division` | `string` | Census division (e.g., `Pacific`, `Mountain`) |
| `createdAt` / `updatedAt` | `string` | Timestamps |

### Category

Top-level scoring bucket. Each category groups related metrics.

| Field | Type | Description |
|-------|------|-------------|
| `id` | `int64` | Primary key |
| `slug` | `string` | URL-safe identifier, UNIQUE |
| `name` | `string` | Display name |
| `description` | `string` | Optional description |
| `defaultWeight` | `float64` | Default weight in overall score (0.0-1.0) |
| `displayOrder` | `int` | Sort order |

### Metric

A measurable indicator within a category. Defines how raw values become scores.

| Field | Type | Description |
|-------|------|-------------|
| `id` | `int64` | Primary key |
| `categoryId` | `int64` | FK to `categories` |
| `slug` | `string` | URL-safe identifier, UNIQUE |
| `name` | `string` | Display name |
| `description` | `string` | Optional description |
| `unit` | `string` | Display unit (`Percent`, `Dollars`, `Per 100k`, etc.) |
| `higherIsBetter` | `bool` | Direction: `true` = higher raw value = better score |
| `normalizationMethod` | `string` | `percentile`, `minmax`, `zscore`, or `fixed` |
| `defaultWeight` | `float64` | Default weight within its category |
| `sourceId` | `*int64` | Optional FK to `data_sources` |
| `active` | `bool` | Whether this metric participates in scoring |
| `createdAt` / `updatedAt` | `string` | Timestamps |

### MetricValue

A single data observation for one state, metric, and year.

| Field | Type | Description |
|-------|------|-------------|
| `id` | `int64` | Primary key |
| `stateId` | `int64` | FK to `states` |
| `metricId` | `int64` | FK to `metrics` |
| `year` | `int` | Data year |
| `value` | `float64` | Numeric observation |
| `sourceRecordId` | `string` | Traceability to upstream source record |
| `importId` | `*int64` | FK to `imports` |
| `createdAt` | `string` | Timestamp |

UNIQUE constraint: `(state_id, metric_id, year, import_id)`.

## Data Provenance

### DataSource

External data origin metadata.

| Field | Type | Description |
|-------|------|-------------|
| `id` | `int64` | Primary key |
| `name` | `string` | Source name |
| `publisher` | `string` | Organization name |
| `sourceUrl` | `string` | URL to the data |
| `license` | `string` | License type |
| `format` | `string` | File format (`csv`, `json`, `pdf`, `html`) |
| `description` | `string` | Free-text description |
| `createdAt` / `updatedAt` | `string` | Timestamps |

### Import

Tracks a data import operation.

| Field | Type | Description |
|-------|------|-------------|
| `id` | `int64` | Primary key |
| `sourceId` | `*int64` | FK to `data_sources` |
| `status` | `string` | `pending`, `running`, `completed`, `completed_with_errors`, `failed` |
| `startedAt` / `completedAt` | `*string` | Timestamps |
| `recordsRead` | `int` | Total rows parsed from CSV |
| `recordsInserted` | `int` | Rows accepted and inserted |
| `recordsRejected` | `int` | Rows rejected by validation |
| `checksum` | `string` | SHA-256 of the uploaded file content |
| `errorSummary` | `string` | Human-readable error summary |

### ImportError

Per-row validation error.

| Field | Type | Description |
|-------|------|-------------|
| `id` | `int64` | Primary key |
| `importId` | `int64` | FK to `imports`, CASCADE delete |
| `rowNumber` | `*int` | 1-indexed row in CSV |
| `fieldName` | `string` | Which column failed |
| `rawValue` | `string` | The offending value |
| `errorMessage` | `string` | Validation error text |

## Application Settings

Key-value store for application-level configuration.

| Field | Type | Description |
|-------|------|-------------|
| `key` | `string` | TEXT PRIMARY KEY |
| `value` | `string` | TEXT value |
| `updatedAt` | `string` | Timestamp |

## Scoring Entities

### ScoringProfile

A weighting configuration for composite scores.

| Field | Type | Description |
|-------|------|-------------|
| `id` | `int64` | Primary key |
| `name` | `string` | Profile name, UNIQUE |
| `description` | `string` | Optional description |
| `isDefault` | `bool` | Whether this is the default profile |
| `isSystem` | `bool` | System profiles cannot be deleted |
| `createdAt` / `updatedAt` | `string` | Timestamps |

### ProfileCategoryWeight

Category weight within a profile.

| Field | Type | Description |
|-------|------|-------------|
| `profileId` | `int64` | FK to `scoring_profiles` |
| `categoryId` | `int64` | FK to `categories` |
| `weight` | `float64` | Weight for this category |

### ProfileMetricWeight

Metric weight within a profile.

| Field | Type | Description |
|-------|------|-------------|
| `profileId` | `int64` | FK to `scoring_profiles` |
| `metricId` | `int64` | FK to `metrics` |
| `weight` | `float64` | Weight for this metric |

### ScoreSnapshot

A cached, calculated composite score.

| Field | Type | Description |
|-------|------|-------------|
| `id` | `int64` | Primary key |
| `profileId` | `int64` | FK to `scoring_profiles` |
| `stateId` | `int64` | FK to `states` |
| `year` | `int` | As-of year |
| `overallScore` | `float64` | 0-100 composite score |
| `completeness` | `float64` | 0.0-1.0 fraction of expected data present |
| `calculatedAt` | `time.Time` | UTC timestamp of calculation |
| `calculationVersion` | `string` | Schema version (currently `1`) |

UNIQUE constraint: `(profile_id, state_id, year, calculation_version)`.

### CategoryScoreSnapshot

Category breakdown for a snapshot.

| Field | Type | Description |
|-------|------|-------------|
| `scoreSnapshotId` | `int64` | FK to `score_snapshots` |
| `categoryId` | `int64` | FK to `categories` |
| `score` | `float64` | 0-100 category score |
| `completeness` | `float64` | 0.0-1.0 completeness within this category |

## Currently Active Metrics

After migration `000010_add_priority_metrics`, 13 metrics across all 5 categories are active:

| Category | Metric Slug | Unit | Direction | Normalization |
|----------|-------------|------|-----------|---------------|
| Economy | `unemployment-rate` | Percent | lowerIsBetter | percentile |
| Economy | `median-household-income` | Dollars | higherIsBetter | percentile |
| Economy | `annual-employment-growth` | Percent | higherIsBetter | percentile |
| Education | `high-school-graduation-rate` | Percent | higherIsBetter | percentile |
| Education | `bachelors-degree-attainment` | Percent | higherIsBetter | percentile |
| Education | `young-adult-college-enrollment` | Percent | higherIsBetter | percentile |
| Health | `life-expectancy` | Years | higherIsBetter | percentile |
| Health | `adult-obesity-prevalence` | Percent | lowerIsBetter | percentile |
| Safety | `violent-crime-rate` | Per 100k | lowerIsBetter | percentile |
| Safety | `traffic-fatalities` | Per 100k | lowerIsBetter | percentile |
| Safety | `property-crime-rate` | Per 100k | lowerIsBetter | percentile |
| Affordability | `cost-of-living-index` | Index (US=100) | lowerIsBetter | percentile |
| Affordability | `renter-housing-cost-burden` | Percent | lowerIsBetter | percentile |

## Inactive Metrics (No Data)

These metrics are defined but deactivated until data is imported:

| Metric Slug | Category | Unit | Direction |
|-------------|----------|------|-----------|
| `uninsured-rate` | Health | Percent | lowerIsBetter |
| `median-rent` | Affordability | Dollars | lowerIsBetter |

# Data Quality Controls

Validation, normalization, and integrity checks throughout the data pipeline.

## CSV Import Validation

The `importer.CSV()` function validates every row before any database writes (all-or-nothing transaction).

### File-Level Checks

| Check | Rule | Error |
|-------|------|-------|
| File size | Max 10 MB | Rejected before parsing |
| Extension | Must be `.csv` | API returns 400 |
| Source existence | `source_id` must exist in `data_sources` | API returns 400 |

### Row-Level Checks

| Field | Validation | Error |
|-------|------------|-------|
| `state_code` | Must be a known two-letter code (uppercased, looked up from `states` table) | Rejected |
| `metric_slug` | Must be a known **active** metric slug (`active=1`) | Rejected |
| `year` | Must parse as integer, between 1900 and current year | Rejected |
| `value` | Must parse as valid `float64` | Rejected |

### Duplicate Detection

Intra-file duplicate: composite key `(state_code, metric_slug, year)` must be unique within a single file.

### Transactional Insert

All valid rows are inserted in a single SQLite transaction. If the transaction fails, no data is committed. This prevents partial imports.

### Error Recording

Each rejected row is recorded in `import_errors` with:
- Row number (1-indexed)
- Field name that failed
- Raw value that was rejected
- Error message

### Import Status Tracking

Import records transition through:

```
pending → running → completed | completed_with_errors | failed
```

A zero-insert import is marked `failed`.

## API-Level Validation

| Check | Implementation |
|-------|---------------|
| File extension | `.csv` required |
| File size | `http.MaxBytesReader` + `io.LimitReader` enforce 10 MB cap |
| Source existence | Database lookup before processing |
| JSON decode limits | Request bodies limited to 1 MB, `DisallowUnknownFields()` rejects unexpected keys |
| Retry logic | Frontend retries on 502/503 up to 3 times with exponential backoff (250ms, 500ms, 750ms) |

## Normalization-Level Quality Checks

| Check | Behavior |
|-------|----------|
| Duplicate detection | `Normalize()` returns error if same `StateID` appears twice |
| Non-finite rejection | `NaN` and `Inf` values cause immediate error |
| Weight validation | `WeightedAverage()` rejects negative weights, `NaN`, and `Inf` weights |
| Score range validation | `WeightedAverage()` rejects scores outside `[0, 100]` |
| Degenerate series | Constant values produce score 50 for all states (no division-by-zero) |

## Completeness Tracking

Every `ScoreSnapshot` and `CategoryScoreSnapshot` carries a `completeness` field (0.0 to 1.0).

```
completeness = sum(included weights) / sum(all positive weights)
```

This measures what fraction of expected positive weight was actually present. A score with completeness < 1.0 is flagged as `Incomplete`, alerting users that the ranking is based on partial data.

## Data Activation Gate

Metrics without seeded data are explicitly set to `active = 0` in migration `000008_scoring_honesty.sql`.

- `MetricRepository.List()` filters by `WHERE active = 1`
- `importer.CSV()` only accepts slugs from `SELECT slug FROM metrics WHERE active=1`

This prevents:
- Empty metrics from diluting category scores via weight redistribution
- Users from importing data for undefined or disabled metrics

## Schema Integrity

| Control | Implementation |
|---------|---------------|
| Foreign keys | `PRAGMA foreign_keys = ON` |
| Journal mode | WAL for concurrent read safety |
| Busy timeout | 5000ms for write contention |
| Connection pool | `MaxOpenConns(1)` appropriate for local SQLite |
| Migrations | Embedded system with `schema_migrations` tracking table |

## Score Snapshot Versioning

`CalculationVersion` (currently `1`) is embedded in the UNIQUE constraint of `score_snapshots`. If the scoring algorithm changes, a new version string produces fresh snapshots without invalidating old ones.

## Import Status Gating

Only imports with status `completed` or `completed_with_errors` are included in the `loadAsOfObservations` SQL query. Failed or in-progress imports never feed into scoring.

## Data Provenance Chain

Every metric value is traceable to its origin:

```
metric_values.import_id → imports (checksum, timestamps, status)
imports.source_id → data_sources (publisher, URL, license)
metric_values.source_record_id → state name from original source
imports.checksum → SHA-256 hash of raw uploaded file
```

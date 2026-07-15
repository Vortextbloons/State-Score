# Scoring Methodology

How raw state data becomes ranked scores.

## Overview

The scoring pipeline has four stages:

```
Raw Values  ──>  Normalization  ──>  Weighted Averages  ──>  Composite Score
(per state)      (0-100)             (category + overall)     (0-100)
```

## Stage 1: Data Loading (As-Of Year)

For each active metric, the system selects the **latest completed import value at or before the requested year** per state. This handles staggered release calendars — for example, combining 2022 life expectancy with 2024 crime data.

SQL logic (from `internal/scoring/recalculate.go`):

```sql
SELECT mv.state_id, mv.value
FROM metric_values mv
JOIN imports i ON i.id = mv.import_id
WHERE mv.metric_id = ?
  AND mv.year <= ?
  AND i.status IN ('completed', 'completed_with_errors')
  AND mv.id = (
    SELECT mv2.id FROM metric_values mv2
    JOIN imports i2 ON i2.id = mv2.import_id
    WHERE mv2.state_id = mv.state_id
      AND mv2.metric_id = mv.metric_id
      AND mv2.year <= ?
      AND i2.status IN ('completed', 'completed_with_errors')
    ORDER BY mv2.year DESC, mv2.id DESC
    LIMIT 1
  )
```

Only imports with status `completed` or `completed_with_errors` are included. Failed or in-progress imports never feed into scoring.

## Stage 2: Normalization

Raw values are converted to 0-100 scores using one of four methods. All scores are clamped to `[0, 100]`. Missing values are omitted entirely (never treated as zero).

### Percentile Rank (`percentile`)

Default for all seeded metrics.

1. Sort observations by value ascending.
2. Ties receive their **average rank** (preventing input order from affecting scores).
3. `score = (rank / (N - 1)) * 100` where rank is 0-indexed average position.
4. If `higherIsBetter == false`: `score = 100 - score`.
5. A single observation gets score 50.

### Min-Max (`minmax`)

1. Compute observed `lo` and `hi` across all present values.
2. `score = ((value - lo) / (hi - lo)) * 100`.
3. If `higherIsBetter == false`: `score = 100 - score`.
4. If all values are equal (`hi == lo`): every state gets 50.
5. Optional `minimum` and `maximum` policy overrides replace observed bounds.

### Fixed Threshold (`fixed`)

1. Uses caller-supplied `minimum` and `maximum` as bounds (required).
2. Same formula as min-max but against policy bounds.
3. Values can exceed bounds; scores are clamped to `[0, 100]`.

### Z-Score (`zscore`)

1. Compute mean and population variance of all present values.
2. Standardize: `z = (value - mean) / stddev`.
3. Map to 0-100 via the standard normal CDF: `score = 50 * (1 + erf(z / sqrt(2)))`.
4. If `higherIsBetter == false`: `score = 100 - score`.
5. If variance is 0 (constant series): every state gets 50.

## Stage 3: Weighted Averages

Metric scores are combined into category scores, then into an overall score.

Formula (from `internal/scoring/scoring.go`):

```
completeness = sum(included weights) / sum(all positive weights)
score = sum(score_i * weight_i) / sum(included weights)
```

**Key behavior**: Missing scores cause their weight to be redistributed to present scores. A category with 0% completeness returns `Incomplete: true` with no score.

### Category Composition

For each category, metric scores are combined using metric weights from the active profile.

### Overall Composition

Category scores are combined using category weights from the active profile.

### Completeness

```
overall_completeness = sum(category_weight * category_completeness) / sum(all positive category_weights)
```

## Default Weights

### Category Weights

| Category | Default Weight |
|----------|----------------|
| Economy | 0.20 |
| Education | 0.20 |
| Health | 0.20 |
| Safety | 0.20 |
| Affordability | 0.20 |

All 5 categories (Economy, Education, Health, Safety, Affordability) currently have active metrics and participate in scoring.

### Metric Weights

Each metric has `default_weight = 0.50` within its category (2 metrics per category, equally weighted).

## Weight Override System

The frontend "Your priorities" page lets users adjust category weights client-side. These re-average the already-normalized backend category scores using the same `WeightedAverage` formula, without re-running normalization.

An optional "normalize to 100%" toggle divides each weight by the total so they always sum to 100%.

## Score Snapshot Versioning

`CalculationVersion` (currently `1`) is embedded in the UNIQUE constraint of `score_snapshots`. If the scoring algorithm changes, a new version string produces fresh snapshots without invalidating old ones. The `ScoreRepository` queries filter by version.

## Worked Example

Given 3 states and 2 metrics in the Safety category:

| State | Violent Crime Rate | Traffic Fatalities |
|-------|-------------------|-------------------|
| CA | 499.0 | 12.8 |
| TX | 438.0 | 15.1 |
| NY | 363.0 | 10.3 |

### Step 1: Percentile Normalization

Violent crime (lowerIsBetter):
- Sorted: NY(363) < TX(438) < CA(499)
- NY: rank 0 → score 100 - (0/2)*100 = 100
- TX: rank 1 → score 100 - (1/2)*100 = 50
- CA: rank 2 → score 100 - (2/2)*100 = 0

Traffic fatalities (lowerIsBetter):
- Sorted: NY(10.3) < CA(12.8) < TX(15.1)
- NY: rank 0 → score 100
- CA: rank 1 → score 50
- TX: rank 2 → score 0

### Step 2: Category Average (equal weights)

| State | Violent Crime | Traffic | Category Score |
|-------|--------------|---------|----------------|
| CA | 0 | 50 | 25 |
| TX | 50 | 0 | 25 |
| NY | 100 | 100 | 100 |

### Step 3: Overall (if Safety weight = 0.20)

Safety contributes 20% of the final score. Health and Affordability also contribute 20% each. The remaining 40% comes from Economy and Education, which now have active metrics.

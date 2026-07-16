# Extending StateScore

## Add a category or metric

1. Add catalog metadata in a new numbered migration; never edit an applied migration. Metrics without data can be hidden by setting `active=0`.
2. Choose a tested normalization method from `internal/scoring`.
3. Import observations using the metric slug.
4. Recalculate snapshots and verify completeness and directionality.
5. Update methodology copy if visible scoring behavior changes.

### Quality gates

If a new metric depends on source data with partial coverage (e.g., FBI property-crime data excludes states below a coverage threshold), add entries to the `metric_value_quality` table after importing values. Set `scoring_eligible=0` with an `exclusion_reason` for observations that should not participate in scoring. See migration `000010_add_priority_metrics.sql` for an example using FBI monthly population coverage.

Catalog-driven screens discover active records automatically. Avoid hard-coded IDs.

The `application_settings` table stores `calculation_version` (currently `"1"`). Increment this value when the scoring algorithm changes.

## Add an API capability

1. Define the path and envelope in `spec/openapi.yaml`.
2. Keep SQL in a repository and orchestration in a feature workflow.
3. Mount the handler under `/api/v1`.
4. Test success, invalid input, and missing data.
5. Add the typed frontend client function and client test.

Prefer bulk endpoints so request counts do not grow with states or metrics.

## Add an import format

Keep parsing and validation in `internal/importer`. Importers must honor cancellation, validate before committing, preserve provenance, and return row-level issues in the common result shape. Imports automatically trigger score recalculation.

## Add background work

Submit application-owned work through `internal/jobs.Manager`. Do not start untracked goroutines or replace the supplied context with `context.Background()`.

## Change persistence

Repositories own SQL. Workflows should consume the smallest storage interface they need. Structural changes belong in `internal/database/migrations`; operational datasets belong in the import workflow. Key-value settings for application state are stored in the `application_settings` table.

# Extending StateScore

## Add a category or metric

1. Add catalog metadata in a new numbered migration; never edit an applied migration.
2. Choose a tested normalization method from `internal/scoring`.
3. Import observations using the metric slug.
4. Recalculate snapshots and verify completeness and directionality.
5. Update methodology copy if visible scoring behavior changes.

Catalog-driven screens discover active records automatically. Avoid hard-coded IDs.

## Add an API capability

1. Define the path and envelope in `spec/openapi.yaml`.
2. Keep SQL in a repository and orchestration in a feature workflow.
3. Mount the handler under `/api/v1`.
4. Test success, invalid input, and missing data.
5. Add the typed frontend client function and client test.

Prefer bulk endpoints so request counts do not grow with states or metrics.

## Add an import format

Keep parsing and validation in `internal/importer`. Importers must honor cancellation, validate before committing, preserve provenance, and return row-level issues in the common result shape.

## Add background work

Submit application-owned work through `internal/jobs.Manager`. Do not start untracked goroutines or replace the supplied context with `context.Background()`.

## Change persistence

Repositories own SQL. Workflows should consume the smallest storage interface they need. Structural changes belong in `internal/database/migrations`; operational datasets belong in the import workflow.

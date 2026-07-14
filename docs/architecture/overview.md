# System Architecture

StateScore is a local-first application: a SvelteKit interface calls a Go JSON API, which persists imported observations and calculated score snapshots in SQLite.

## Components and dependency direction

```text
Svelte routes -> frontend API client -> /api/v1
                                      |
Go HTTP handlers -> domain workflows -> repositories -> SQLite
                                      |
                              managed background jobs
```

- `internal/app` is the composition root and owns configuration, database, HTTP, jobs, and shutdown.
- `internal/api` validates HTTP input and emits the versioned JSON contract.
- `internal/scoring` contains normalization and score calculation behavior.
- `internal/importer` validates CSV observations and writes imports atomically.
- `internal/repositories` owns persistent queries.
- `internal/jobs` owns cancellation, accounting, and shutdown of background work.
- `frontend/src/lib/api` is the typed browser boundary.

## Workflows

Rankings load catalogs, all observations through bulk `GET /values`, and a canonical score snapshot. Imports are submitted to the managed job lifecycle, validated before insertion, and trigger snapshot recalculation. Scoring selects the latest observation at or before the requested year, normalizes active metrics, applies profile weights, and persists versioned snapshots.

## Contracts and boundaries

Browser input and files are untrusted. `spec/openapi.yaml` is the route-level API contract. SQL migrations are embedded and applied in lexical order. Imported datasets are data, not schema; future large datasets should use imports rather than seed migrations.

See [Extending StateScore](../guides/extending.md) for supported extension workflows.

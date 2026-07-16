# System Architecture

StateScore is a local-first application: a SvelteKit 2 static SPA (Svelte 5 runes, built via @sveltejs/adapter-static with 200.html fallback) calls a Go JSON API served on 127.0.0.1. The Go backend persists imported observations and calculated score snapshots in SQLite via modernc.org/sqlite (pure Go, no CGO).

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
- `web/embed.go` embeds the compiled frontend (`web/dist/`) via `//go:embed all:dist`
- `internal/webui/handler.go` serves the embedded SPA
- `internal/security/` and `internal/metrics/` exist as empty packages (no source files yet)

## Workflows

Rankings load catalogs, all observations through bulk `GET /values`, and a canonical score snapshot. Imports are submitted to the managed job lifecycle, validated before insertion, and trigger snapshot recalculation. Scoring selects the latest observation at or before the requested year, normalizes active metrics, applies profile weights, and persists versioned snapshots. Observations with `scoring_eligible = 0` in the `metric_value_quality` table are excluded during the as-of observation query (`loadAsOfObservations` in `internal/scoring/recalculate.go`).

## Contracts and boundaries

Browser input and files are untrusted. Routes use Go 1.22+ `http.ServeMux` with method-pattern matching. The server listens on `127.0.0.1` (default port 8787, with automatic fallback up to +49 if occupied). No authentication, TLS, or CORS headers — localhost-only design. SQL migrations are embedded and applied in lexical order. Imported datasets are data, not schema; future large datasets should use imports rather than seed migrations.

Graceful shutdown: SIGINT/SIGTERM stops the HTTP server first, then cancels background jobs (10s timeout).

## Docker deployment

A multi-stage Docker build (`Dockerfile`) compiles the Go binary and embeds the frontend build in a single stage, then copies the binary into an `alpine:3.21` runtime image. The container runs on `0.0.0.0:8787`, disables automatic browser opening (`STATESCORE_NO_BROWSER=1`), and stores data in a mounted volume at `/data`. The `docker-compose.yml` provides a single-service setup with a named volume for persistence and `restart: unless-stopped`.

> **Note:** The OpenAPI spec at `spec/openapi.yaml` is significantly outdated — many routes are missing.

See [Extending StateScore](../guides/extending.md) for supported extension workflows.

# Components

Overview of project components and their responsibilities.

## Frontend

SvelteKit 2 + Svelte 5 (runes mode) static SPA, written in TypeScript. Built with `@sveltejs/adapter-static`, output to `web/dist/`. Includes 7 routes (Overview (`/`), Rankings (`/rankings`), Compare (`/compare`), Your Priorities (`/scoring`), Data Workshop (`/data`), Methodology (`/methodology`), State Profile (`/states/[code]`)), API client with retry logic in `src/lib/api/client.ts`, and data orchestration in `src/lib/scores.ts`. No SSR (`ssr = false`). Uses `localStorage` for saved scoring perspectives.

## Backend

Go HTTP server (`cmd/statescore` + `internal/`) using Go 1.22+ `http.ServeMux` with method-pattern routing. SQLite via `modernc.org/sqlite` (pure Go, no CGO). Repository pattern with dependency injection interface stores. Background job manager for async CSV imports. Embedded SQL migrations and frontend build via `//go:embed`. Default port: 8787.

Migration `000010_add_priority_metrics.sql` added five priority metrics (one per category: annual-employment-growth, young-adult-college-enrollment, adult-obesity-prevalence, property-crime-rate, renter-housing-cost-burden) with bundled 2024 observations. It also introduced the `metric_value_quality` table, which gates metric values from scoring via the `scoring_eligible` column (checked in `loadAsOfObservations`).

Migration `000011_add_foundational_metrics.sql` adds labor-force participation, a four-assessment NAEP achievement composite, uninsured rate, age-adjusted homicide mortality, and owner housing-cost burden with complete official 2024 observations. Catalog and methodology pages render these active metrics without metric-specific frontend code.

## Connection

In development, the Vite dev server proxies `/api` requests to `http://127.0.0.1:8080`. In production, the Go server embeds the static frontend build and serves it directly.

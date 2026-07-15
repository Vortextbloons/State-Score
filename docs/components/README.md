# Components

Overview of project components and their responsibilities.

## Frontend

SvelteKit 2 + Svelte 5 (runes mode) static SPA, written in TypeScript. Built with `@sveltejs/adapter-static`, output to `web/dist/`. Includes 7 routes (Home, Rankings, Compare, Your Priorities, State Profile, Data Workshop, Methodology), API client with retry logic in `src/lib/api/client.ts`, and data orchestration in `src/lib/scores.ts`. No SSR (`ssr = false`). Uses `localStorage` for saved scoring perspectives.

## Backend

Go HTTP server (`cmd/statescore` + `internal/`) using Go 1.22+ `http.ServeMux` with method-pattern routing. SQLite via `modernc.org/sqlite` (pure Go, no CGO). Repository pattern with dependency injection interface stores. Background job manager for async CSV imports. Embedded SQL migrations and frontend build via `//go:embed`. Default port: 8787.

## Connection

In development, the Vite dev server proxies `/api` requests to `http://127.0.0.1:8080`. In production, the Go server embeds the static frontend build and serves it directly.

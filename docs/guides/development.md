# Development Guide

Install Go, Node.js, and npm, then install both dependency sets:

```bash
npm install
npm install --prefix frontend
```

Run canonical commands from the repository root:

```bash
npm run dev              # frontend + backend concurrently
npm run dev:frontend     # SvelteKit only on :5173
npm run dev:backend      # Go server on :8080
npm run check            # svelte-kit sync + svelte-check
npm test                 # Vitest (frontend tests)
npm run lint             # Prettier + ESLint
npm run build            # builds frontend
npm run combine-docs     # rebuilds docs/ALL.md
```

Make targets are also available: `make dev`, `make dev-backend`, `make dev-backend-live`, `make build`, `make frontend-build`, `make go-build`, `make test`, `make tidy`.

The frontend is SvelteKit 2 with Svelte 5 runes mode and TypeScript. The Vite dev server runs on `127.0.0.1:5173` and proxies `/api` to `http://127.0.0.1:8080`. The backend uses Go's standard HTTP library and SQLite; the default port is 8787 (dev overrides to 8080 via `STATESCORE_PORT`).

For Go hot-reload during development, install [Air](https://github.com/air-verse/air):

```bash
go install github.com/air-verse/air@latest
```

Configuration is in `.air.toml`.

Run Go tests with `go test ./cmd/... ./internal/... ./web/...`. The root `tests/` directory is empty (placeholder for future integration tests).

For API changes, update `spec/openapi.yaml` and frontend types. For schema changes, add a migration under `internal/database/migrations`. After documentation changes, run `npm run combine-docs`; never edit `docs/ALL.md` manually.

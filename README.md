# State Score

Rank and compare U.S. states across quality-of-life categories.

## Features

- **Side-by-side comparison** — Pick any two states and compare scores across multiple dimensions.
- **Your Priorities** — Adjust category weights to score states by what matters to you.
- **CSV Data Import** — Workshop tool for importing and previewing new datasets.
- **Methodology** — Scoring, normalization, and weighting are documented in full.
- **Local-first** — Entirely self-contained; no external servers or telemetry.

## Stack

| Layer   | Technology                      |
|---------|---------------------------------|
| Backend | Go HTTP server (`cmd/statescore`) |
| Frontend| SvelteKit + TypeScript (`frontend/`) |
| Database| SQLite with versioned migrations |

The frontend is embedded into the Go binary at build time, producing a single deployable artifact.

## Quick Start

```bash
# install dependencies
npm install
npm install --prefix frontend

# run backend + frontend concurrently
npm run dev
```

- Frontend: [http://localhost:5173](http://localhost:5173)
- Backend API: [http://localhost:8080](http://localhost:8080)

See [docs/guides/development.md](docs/guides/development.md) for detailed setup instructions.

## Documentation

| File | Description |
|------|-------------|
| [docs/INDEX.md](docs/INDEX.md) | Full documentation index |
| [docs/guides/development.md](docs/guides/development.md) | Setup and development guide |
| [docs/reference/api.md](docs/reference/api.md) | API reference |
| [docs/reference/configuration.md](docs/reference/configuration.md) | Configuration reference |

# Development Guide

Install Go, Node.js, and npm, then install both dependency sets:

```bash
npm install
npm install --prefix frontend
```

Run canonical commands from the repository root:

```bash
npm run dev
npm run check
npm test
npm run lint
npm run build
npm run combine-docs
```

The frontend is SvelteKit with TypeScript. The backend uses Go's standard HTTP library and SQLite.

For API changes, update `spec/openapi.yaml` and frontend types. For schema changes, add a migration under `internal/database/migrations`. After documentation changes, run `npm run combine-docs`; never edit `docs/ALL.md` manually.

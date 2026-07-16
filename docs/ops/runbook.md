# Operations Runbook

Runbook and troubleshooting procedures for the State Score application.

## Building

```bash
make build          # builds frontend then Go binary
# or step-by-step:
npm run build       # frontend -> web/dist/
go build -o bin/statescore ./cmd/statescore
# Windows produces bin/statescore.exe
```

The Go binary embeds the frontend (`web/dist/`) and SQL migrations at compile time. No external runtime dependencies — pure Go SQLite via `modernc.org/sqlite` (no CGO).

## Running

```bash
./bin/statescore
```

- **Port**: `8787` by default, configurable via `STATESCORE_PORT`
- **Host**: `127.0.0.1` by default, configurable via `STATESCORE_HOST`
- **Database**: auto-created on first run at the OS-specific data directory
- **Browser**: auto-opens on startup; disable with `STATESCORE_NO_BROWSER=1`
- **Shutdown**: graceful on SIGINT/SIGTERM with 10-second timeout

### Port Fallback

If `STATESCORE_PORT` is busy, the app scans up to 49 subsequent ports (`port+1` through `port+49`) and uses the first available. Useful if an old instance lingers or another service occupies the configured port.

### Data Directory

| Platform | Path |
|---|---|
| Windows | `%LOCALAPPDATA%\StateScore` |
| macOS   | `~/Library/Application Support/StateScore` |
| Linux   | `$XDG_DATA_HOME/statescore` or `~/.local/share/statescore` |

## Security Boundaries

| Concern | Status |
|---|---|
| Authentication | None — no auth middleware, no passwords, no API keys |
| TLS | None — plain HTTP (localhost-only, no CORS needed) |
| CORS | Not set — Vite dev proxy handles cross-origin in dev; production is localhost-only (Docker sets `STATESCORE_HOST=0.0.0.0`, making the server accessible from any network interface — bind to a specific IP or add a reverse proxy if network exposure is a concern) |
| SQL injection | Mitigated — parameterized queries used throughout |
| CSV upload limit | 10 MB max via `http.MaxBytesReader` + `io.LimitReader`; `.csv` extension required |
| JSON request limit | 1 MB; `DisallowUnknownFields()` rejects unexpected keys |
| Static file serving | Only GET/HEAD allowed; paths under `/api/` are blocked |
| `internal/security/` | Package directory exists but is **empty** — no security code is currently enforced here |

## Docker Deployment

```bash
docker compose build
docker compose up -d        # starts on host port :8787
docker compose logs -f      # tail server logs
docker compose down         # stop container
```

Alternatively, build with plain Docker:

```bash
docker build -t statescore .
docker run -d -p 8787:8787 -v statescore-data:/data statescore
```

The image runs the same Go binary in an Alpine container. Key differences from a local run:

- **Host**: `STATESCORE_HOST=0.0.0.0` — the server listens on all network interfaces (not just localhost).
- **Browser**: `STATESCORE_NO_BROWSER=1` — no browser auto-open inside the container.
- **Data**: `XDG_DATA_HOME=/data` — persistent SQLite database stored on a Docker volume (`statescore-data`).
- **Restart**: `restart: unless-stopped` in the compose file keeps the container running after crashes or host reboots.
- **Logs**: Use `docker compose logs -f` instead of looking for stdout; there is no persistent log file.

## Common Issues

### Server won't start / Port already in use

- The app auto-scans up to `port+49` if the configured port is busy — check logs for which port it bound to.
- Set `STATESCORE_PORT` to a different value if another instance is already running.
- Kill stale processes if needed.

### "No frontend assets" / 503 on frontend routes

`web/dist/` is empty or missing. The Go binary embeds the frontend at build time — if the frontend wasn't built first, assets won't be embedded.

**Fix**: Run `make build` (or `npm run build && go build -o bin/statescore ./cmd/statescore`) to build frontend assets before the Go binary.

### Database issues

- **WAL mode + 5s busy timeout** handles most contention.
- **`MaxOpenConns(1)`** ensures a single writer — locked-database errors in normal operation are unexpected.
- Check file permissions on the data directory if the database won't open.

### Migration failures

Migrations run automatically on startup from `internal/database/migrations/`. If a migration fails:

1. The app will log the SQL error and exit.
2. Verify the migration SQL is valid for SQLite (modernc.org/sqlite).
3. Check that the database file isn't from an incompatible version.

### Browser doesn't open (macOS/Linux)

- macOS requires `open` in PATH.
- Linux requires `xdg-open` in PATH.
- Set `STATESCORE_NO_BROWSER=1` if you prefer to open the URL manually.

### CORS errors in development

Vite's dev proxy (`npm run dev:frontend`) handles CORS when proxying API calls to the Go backend. If you access the Go server directly from a browser on a different origin, there are no CORS headers — this is expected for a localhost-only app.

# Configuration

Environment variables and derived configuration values.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `STATESCORE_PORT` | `8787` | HTTP listen port. Dev backend overrides to 8080. Falls back through ports up to +49 if busy. |
| `STATESCORE_NO_BROWSER` | *(unset)* | Set to `1` to skip opening browser on startup. |

## Derived Configuration

| Field | Value |
|-------|-------|
| `Host` | `127.0.0.1` |
| `DataDir` | `%LOCALAPPDATA%/StateScore` (Windows), `~/Library/Application Support/StateScore` (macOS), `$XDG_DATA_HOME/statescore` or `~/.local/share/statescore` (Linux) |
| `DatabasePath` | `<DataDir>/statescore.db` |
| `OpenBrowser` | `true` unless `STATESCORE_NO_BROWSER=1` |
| `Version` | `0.1.0` (hardcoded) |

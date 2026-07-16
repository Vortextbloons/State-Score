# API Reference

Base path: `/api/v1`. Error responses use `{"error":{"code":"str","message":"str","details":obj}}`.
The [OpenAPI spec](../../spec/openapi.yaml) is incomplete and may be outdated.

## Operations

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Returns `{"status":"ok","version":"0.1.0","startedAt":"...","app":"StateScore"}` |
| GET | `/status` | Returns `{"status":"ready\|degraded","version":"0.1.0","databaseReady":bool,"activeImports":int,"startedAt":"..."}` |

## Catalogs

| Method | Path | Description |
|--------|------|-------------|
| GET | `/states?region=` | List states with population estimate and year, optional region filter |
| GET | `/states/{code}` | Get state by 2-letter code (e.g., CA), including population estimate and year |
| GET | `/regions` | List all regions |
| GET | `/categories` | List scoring categories |
| GET | `/metrics?category_id=` | List metrics, optional category filter |
| GET | `/metrics/{metricId}` | Get single metric |
| GET | `/profiles` | List scoring profiles |
| GET | `/profiles/{profileId}` | Get profile with category weights |
| GET | `/profiles/default` | Get default profile with weights |

## Values

`GET /values?state_id=&year=` — Get metric values with optional filters. Each value includes a `quality` object (`reportingCoverage`, `participatingAgencies`, `populationCovered`, `dataRevision`, `scoringEligible`, `exclusionReason`) when source-specific metadata is available.

## Sources & Imports

| Method | Path | Description |
|--------|------|-------------|
| GET | `/sources` | List data sources |
| POST | `/sources` | Create source (JSON body; 422 if name missing or format not CSV) |
| PUT | `/sources/{sourceId}` | Update source |
| GET | `/imports?limit=50` | List imports |
| POST | `/imports` | Upload CSV (multipart: `source_id` int64 + `file` .csv, max 10MB). Returns 202. |
| GET | `/imports/{importId}` | Get import with validation errors |
| GET | `/public-sources` | List registered refresh adapters and availability |
| POST | `/public-sources/refresh` | Queue one or more official-source refreshes |

CSV columns: `state_code`, `metric_slug`, `year`, `value` (required); `source_record_id` (optional).

Refresh request: `{"adapterIds":["census-college-enrollment","cdc-adult-obesity"],"year":2024}`. An empty adapter list refreshes all available adapters. The response maps adapter IDs to import IDs; progress appears in the import ledger.

## Scores

| Method | Path | Description |
|--------|------|-------------|
| GET | `/scores?profile_id=&year=` | Scoreboard — returns `{profileId, year, asOfYear, calculationVersion, scores: [{stateId, overallScore, completeness, calculatedAt, calculationVersion, categories}]}` |
| POST | `/scores/recalculate` | Body `{"profileId":int64,"year":int}` — returns `{profileId, year, statesCalculated:int}` |

## Frontend

All other paths (`/`, etc.) serve the embedded SvelteKit SPA. Returns 503 if no frontend build is available.

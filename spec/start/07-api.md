## 13. API specification

Base path:

```text
/api/v1
```

### Application status

```http
GET /api/v1/status
```

Example response:

```json
{
  "status": "ready",
  "version": "0.1.0",
  "databaseReady": true,
  "activeImports": 0
}
```

### List states

```http
GET /api/v1/states
```

Optional query parameters:

```text
year
profile
sort
direction
region
search
```

### Read state

```http
GET /api/v1/states/{code}
```

### Compare states

```http
GET /api/v1/compare?states=UT,CO,ID&year=2025&profile=balanced
```

### Rankings

```http
GET /api/v1/rankings
```

Parameters:

```text
year
profile
category
sort
direction
```

### Categories

```http
GET /api/v1/categories
```

### Metrics

```http
GET /api/v1/metrics
GET /api/v1/metrics/{metricId}
```

### Scoring profiles

```http
GET    /api/v1/profiles
POST   /api/v1/profiles
GET    /api/v1/profiles/{profileId}
PUT    /api/v1/profiles/{profileId}
DELETE /api/v1/profiles/{profileId}
```

### Imports

```http
GET  /api/v1/imports
POST /api/v1/imports
GET  /api/v1/imports/{importId}
POST /api/v1/imports/{importId}/cancel
```

### Sources

```http
GET  /api/v1/sources
POST /api/v1/sources/{sourceId}/download
```

### Recalculate scores

```http
POST /api/v1/scores/recalculate
```

### Export

```http
GET /api/v1/export/rankings?format=csv
GET /api/v1/export/comparison?states=UT,CO&format=json
```

### Shutdown

```http
POST /api/v1/application/shutdown
```

The shutdown endpoint must require a local session token and explicit confirmation.

---

## 14. API response conventions

Successful response:

```json
{
  "data": {},
  "meta": {}
}
```

Error response:

```json
{
  "error": {
    "code": "INVALID_STATE_CODE",
    "message": "The state code must be a valid two-letter U.S. state code.",
    "details": {}
  }
}
```

Expected HTTP statuses:

```text
200 Successful request
201 Resource created
202 Background operation accepted
204 Successful request with no response body
400 Invalid request
404 Resource not found
409 Conflict
422 Validation failed
500 Internal application error
503 Application temporarily unavailable
```

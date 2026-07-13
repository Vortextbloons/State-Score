# StateScore — Product and Technical Specification

## 1. Product summary

**StateScore** is a locally installed web application for comparing U.S. states across categories such as affordability, economy, education, health, safety, and quality of life.

When the user starts the application:

1. A Go process starts a local HTTP server.
2. The server loads the local SQLite database.
3. The application selects an available localhost port.
4. The user’s default browser opens automatically.
5. The browser displays the Svelte interface.
6. All application data remains on the user’s computer unless the user explicitly downloads updated public datasets.

Example address:

```text
http://127.0.0.1:8787
```

The application must not require an internet connection after its datasets have been downloaded and imported.

---

## 2. Primary goals

The application must allow users to:

* View an overall score for every U.S. state.
* Rank states from highest to lowest.
* Compare two or more states.
* inspect category and metric scores.
* Change category weights.
* Recalculate rankings immediately.
* View the source and year of every metric.
* Import updated datasets.
* Use the application locally without creating an account.
* Export comparisons and rankings.

The score must be transparent. The application must never present a state’s overall score without showing how it was calculated.

---

## 3. Non-goals for the first version

The first release will not include:

* User accounts.
* Cloud synchronization.
* Social features.
* Mobile-native applications.
* Automatic scraping of arbitrary websites.
* Predictions about future state performance.
* AI-generated recommendations.
* County- or city-level comparisons.
* Collaborative editing.
* Public internet hosting.

These features may be considered after the local desktop-style application is stable.

---

## 4. Technology stack

### Frontend

* SvelteKit
* TypeScript
* Static SvelteKit build
* Standard CSS or a small component library
* Charting library for comparison and trend charts

SvelteKit will use its static adapter so the frontend can be built into static HTML, JavaScript, and CSS files.

### Backend

* Go
* Go `net/http`
* JSON REST API
* Background data-import jobs
* Local scoring engine

Go’s standard HTTP package will serve both the application frontend and the JSON API.

### Storage

* SQLite
* Local database file
* SQL migrations
* Optional SQLite WAL mode
* CSV and JSON import support

### Distribution

The production build should contain:

* One Go executable.
* Embedded frontend files.
* An automatically created SQLite database.
* Optional starter datasets.

The built Svelte files can be compiled into the Go executable with Go’s `embed` package.

The project will not use Electron, Wails, Tauri, or another embedded-webview framework.

---

## 5. Runtime architecture

```text
Default browser
      |
      | HTTP requests
      v
http://127.0.0.1:<port>
      |
      v
Go application
  ├── Static Svelte frontend
  ├── REST API
  ├── Scoring engine
  ├── Dataset importer
  └── SQLite database
```

The frontend and API must use the same origin:

```text
Frontend: http://127.0.0.1:8787/
API:      http://127.0.0.1:8787/api/v1/
```

This avoids unnecessary cross-origin configuration.

---

## 6. Application startup behavior

When the executable starts, it must:

1. Determine the operating-system-specific application data directory.
2. Create the application directory when it does not exist.
3. Open or create the SQLite database.
4. Run pending database migrations.
5. Load configuration.
6. Attempt to bind to `127.0.0.1:8787`.
7. Use another available local port if port 8787 is occupied.
8. Start the HTTP server.
9. Wait until the server is ready.
10. Open the application URL in the default browser.
11. Continue running until the user stops the process or selects **Shut Down Application**.

The server must bind to `127.0.0.1`, not `0.0.0.0`, unless the user deliberately enables network access in a future advanced setting.

Example startup output:

```text
StateScore is running.

Open: http://127.0.0.1:8787
Data: /home/user/.local/share/statescore/statescore.db

Press Ctrl+C to stop.
```

If the browser cannot be opened automatically, the terminal must display a clickable or copyable local address.

Starting the executable a second time should either:

* Open the existing application URL, or
* Detect that another StateScore process is already using the configured port and exit cleanly.

---

## 7. Shutdown behavior

The application must support:

* `Ctrl+C` shutdown from the terminal.
* Operating-system termination signals.
* A **Shut Down Application** button in the settings page.
* Graceful completion or cancellation of active imports.
* Closing the SQLite connection.
* Graceful HTTP server shutdown.

Closing the browser tab does not automatically stop the Go process.

The interface should clearly communicate this:

```text
Closing this browser tab does not shut down StateScore.
Use Settings → Shut Down Application or stop the terminal process.
```

---

## 8. Main application pages

### 8.1 Home dashboard

The dashboard displays:

* Highest-ranked states.
* Lowest-ranked states.
* Current scoring profile.
* Number of metrics included.
* Dataset year range.
* Latest import date.
* Warnings about missing or outdated data.
* Quick state-comparison controls.

Example:

```text
Current profile: Balanced
States ranked: 50
Metrics used: 18
Data coverage: 2024–2026
Last updated: July 12, 2026
```

### 8.2 Rankings page

The rankings page displays a sortable table:

| Rank | State         | Overall | Economy | Education | Health | Safety | Affordability |
| ---: | ------------- | ------: | ------: | --------: | -----: | -----: | ------------: |
|    1 | Example State |    86.4 |    91.0 |      84.2 |   88.5 |   82.1 |          79.3 |

Required controls:

* Sort by overall score.
* Sort by category.
* Search by state name.
* Filter by region.
* Select scoring profile.
* Change ranking year.
* Add a state to comparison.
* Export the current table.
* Open a state profile.

The table must show when a score is based on incomplete data.

### 8.3 State profile page

Route:

```text
/states/UT
```

The page displays:

* State name and abbreviation.
* Overall score.
* Overall rank.
* Category scores.
* Category ranks.
* Metric values.
* Metric scores.
* Data years.
* Data sources.
* Missing-data warnings.
* Historical trends when multiple years are available.
* An explanation of score calculation.

Example:

```text
Utah

Overall score: 84.9
Overall rank:   4 of 50

Economy:       91.2
Education:     82.1
Health:        80.4
Safety:        86.0
Affordability: 75.3
```

### 8.4 State comparison page

Route example:

```text
/compare?states=UT,CO,ID
```

The comparison page must support between two and five states.

It displays:

* Overall scores.
* Category scores.
* Individual metric values.
* State ranks.
* Bar charts.
* Radar chart as an optional view.
* Historical line charts.
* Source and year differences.
* Clear indication of which metrics are better when higher or lower.

Users must be able to:

* Add or remove states.
* Change category weights.
* Hide categories.
* Export results.
* Copy a comparison summary.

### 8.5 Custom scoring page

Users can define how much each category matters.

Example:

```text
Economy        25%
Education      20%
Health         20%
Safety         20%
Affordability  15%
                ---
                100%
```

Requirements:

* Weight sliders and numeric inputs.
* Weights must total 100%.
* Automatic normalization option.
* Reset to default.
* Save profile.
* Duplicate profile.
* Rename profile.
* Delete custom profile.
* Show ranking changes immediately.

Example profiles:

* Balanced
* Best for Families
* Most Affordable
* Career Focused
* Health and Safety
* Retirement
* Custom

Preset names must be presented as configurable scoring perspectives, not objective declarations.

### 8.6 Data sources page

This page displays:

* Dataset name.
* Publishing organization.
* Source address.
* Import date.
* Data year.
* Number of records.
* File format.
* License or usage notes.
* Metrics supplied by the dataset.
* Import status.
* Validation warnings.

### 8.7 Data import page

The first version must support manual imports of:

* CSV files.
* JSON files.
* Preconfigured remote datasets.

The import workflow:

1. Select a source definition.
2. Select a local file or download the configured source.
3. Preview detected columns.
4. Validate state identifiers.
5. Map source columns to metrics.
6. Preview transformed values.
7. Run the import.
8. Display errors and warnings.
9. Recalculate affected scores.
10. Save an import report.

Imports must run in a background goroutine so the HTTP server remains responsive.

### 8.8 Methodology page

This page must explain:

* Every category.
* Every metric.
* Metric direction.
* Normalization method.
* Missing-data handling.
* Weighting method.
* Outlier handling.
* State exclusions.
* Data-year rules.
* Overall-score formula.
* Limitations of the ranking.

### 8.9 Settings page

Settings include:

* Preferred localhost port.
* Automatically open browser.
* Automatically check for dataset updates.
* Default scoring profile.
* Default comparison year.
* Export directory.
* Database backup.
* Database restore.
* Reset local data.
* Shut down application.

---

## 9. Scoring categories

The initial version should use five categories.

### Economy

Possible metrics:

* Unemployment rate.
* Median household income.
* Poverty rate.
* Employment growth.
* Gross state product per capita.

### Education

Possible metrics:

* High-school graduation rate.
* Bachelor’s degree attainment.
* Student-to-teacher ratio.
* Standardized educational outcomes.
* Education spending per student.

### Health

Possible metrics:

* Life expectancy.
* Uninsured rate.
* Preventable hospitalization rate.
* Obesity rate.
* Access to primary-care providers.

### Safety

Possible metrics:

* Violent crime rate.
* Property crime rate.
* Traffic fatalities per 100,000 people.
* Workplace fatality rate.
* Emergency-response indicators.

### Affordability

Possible metrics:

* Median rent.
* Median home-price-to-income ratio.
* Cost-of-living index.
* State and local tax burden.
* Utility costs.

The MVP should begin with approximately two metrics per category. Additional metrics should only be added after their definitions and coverage are verified.

---

## 10. Metric definition requirements

Each metric must contain:

```text
Name
Slug
Description
Category
Unit
Source
Data year
Higher-is-better flag
Normalization method
Weight within category
Minimum accepted value
Maximum accepted value
Missing-data rule
```

Example:

```text
Name: Unemployment rate
Slug: unemployment-rate
Category: Economy
Unit: Percent
Higher is better: No
Normalization: Percentile
Weight within category: 50%
```

Metrics must use rates, percentages, per-capita measures, or another normalized unit whenever raw totals would unfairly favor states with larger or smaller populations.

---

## 11. Score calculation

### 11.1 Metric normalization

Every metric is converted to a score from 0 to 100.

For a higher-is-better metric:

```text
score = ((value - minimum) / (maximum - minimum)) × 100
```

For a lower-is-better metric:

```text
score = ((maximum - value) / (maximum - minimum)) × 100
```

The implementation should support multiple normalization methods:

* Min-max.
* Percentile rank.
* Z-score converted to a 0–100 scale.
* Fixed policy thresholds.

Percentile scoring is recommended as the default because basic min-max scoring can be distorted by extreme values.

### 11.2 Category score

```text
category score =
sum(metric score × metric weight)
/
sum(included metric weights)
```

### 11.3 Overall score

```text
overall score =
economy score × economy weight
+ education score × education weight
+ health score × health weight
+ safety score × safety weight
+ affordability score × affordability weight
```

### 11.4 Missing values

The application must never silently convert missing data to zero.

Supported missing-data behavior:

* Exclude the metric and redistribute its weight.
* Mark the category as incomplete.
* Exclude the state from that ranking.
* Use an explicitly configured imputation method.

The default should exclude the missing metric, redistribute its weight within the category, and show a visible completeness warning.

### 11.5 Score reproducibility

Every stored score calculation must reference:

* Scoring profile.
* Profile version.
* Metric version.
* Dataset import.
* Data year.
* Calculation timestamp.

A user must be able to understand why an old exported ranking differs from a newly calculated ranking.

---

## 12. Database design

### `states`

```text
id
code
name
region
division
created_at
updated_at
```

### `categories`

```text
id
slug
name
description
default_weight
display_order
```

### `metrics`

```text
id
category_id
slug
name
description
unit
higher_is_better
normalization_method
default_weight
source_id
active
created_at
updated_at
```

### `metric_values`

```text
id
state_id
metric_id
year
value
source_record_id
import_id
created_at
```

Unique constraint:

```text
state_id + metric_id + year + import_id
```

### `data_sources`

```text
id
name
publisher
source_url
license
format
description
created_at
updated_at
```

### `imports`

```text
id
source_id
status
started_at
completed_at
records_read
records_inserted
records_rejected
checksum
error_summary
```

### `import_errors`

```text
id
import_id
row_number
field_name
raw_value
error_message
```

### `scoring_profiles`

```text
id
name
description
is_default
is_system
created_at
updated_at
```

### `profile_category_weights`

```text
profile_id
category_id
weight
```

### `profile_metric_weights`

```text
profile_id
metric_id
weight
```

### `score_snapshots`

```text
id
profile_id
state_id
year
overall_score
completeness
calculated_at
calculation_version
```

### `category_score_snapshots`

```text
score_snapshot_id
category_id
score
completeness
```

### `application_settings`

```text
key
value
updated_at
```

---

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

---

## 15. Frontend state and behavior

The Svelte frontend must:

* Use TypeScript.
* Treat the Go API as the source of truth.
* Display loading states.
* Display recoverable error messages.
* Preserve custom weights during navigation.
* Debounce weight-slider recalculations.
* Use URL query parameters for shareable local comparison states.
* Support keyboard navigation.
* Avoid requiring JavaScript-generated content for initial error pages.
* Display data-year and source information near important scores.

Changing weights should update the visible ranking immediately. The frontend may calculate temporary preview rankings, but the Go backend must produce the canonical saved result.

---

## 16. Localhost security requirements

Although the application runs locally, it must not assume that every browser request is trusted.

Required protections:

* Bind only to `127.0.0.1` and optionally `::1`.
* Reject unexpected `Host` headers.
* Do not enable broad CORS access.
* Use same-origin API requests.
* Generate a random session token at startup.
* Require the token for write, import, reset, and shutdown operations.
* Use CSRF protection for state-changing requests.
* Set appropriate security headers.
* Limit uploaded file sizes.
* Validate CSV and JSON input.
* Never execute imported content.
* Escape values rendered in the browser.
* Use parameterized SQL queries.
* Store no secret credentials in frontend files.
* Reject paths containing directory traversal sequences.
* Do not serve arbitrary files from the user’s computer.

The browser-opening URL may contain a short-lived bootstrap token. The frontend should exchange it for a session cookie and then remove the token from the visible address.

---

## 17. Data validation requirements

Imported data must be checked for:

* Valid state name or code.
* Recognized year.
* Numeric metric value.
* Expected unit.
* Duplicate records.
* Impossible values.
* Missing required fields.
* Unexpected state or territory identifiers.
* Data outside configured bounds.
* Dataset schema changes.

Example validation rules:

```text
Unemployment rate must be between 0 and 100.
State code must match a known state.
Year must be between 1900 and the current year.
A state cannot have two active values for the same metric and year.
```

Rejected records must be included in the import report rather than silently discarded.

---

## 18. Export requirements

Supported exports:

* Ranking CSV.
* Ranking JSON.
* State-comparison CSV.
* State-comparison JSON.
* Printable HTML report.
* Methodology summary.

Every export must include:

* Export timestamp.
* Application version.
* Scoring-profile name.
* Category weights.
* Metric weights.
* Data years.
* Data sources.
* Missing-data warnings.
* Calculation version.

---

## 19. Development environment

### Development mode

Run two development servers:

```text
Svelte development server: http://127.0.0.1:5173
Go API server:             http://127.0.0.1:8080
```

The Svelte development server should proxy `/api` requests to the Go server.

Example workflow:

```bash
# Terminal 1
cd frontend
npm run dev

# Terminal 2
cd backend
go run ./cmd/statescore
```

### Production mode

1. Build the Svelte static frontend.
2. Copy or generate the frontend output inside the Go embedding directory.
3. Compile the Go executable.
4. Run the executable.
5. Serve the embedded frontend and API from one localhost origin.

Example:

```bash
npm run build
go build -o statescore ./cmd/statescore
./statescore
```

---

## 20. Suggested repository structure

```text
statescore/
├── cmd/
│   └── statescore/
│       └── main.go
├── internal/
│   ├── api/
│   ├── app/
│   ├── browser/
│   ├── config/
│   ├── database/
│   ├── importer/
│   ├── metrics/
│   ├── scoring/
│   ├── security/
│   └── shutdown/
├── migrations/
├── frontend/
│   ├── src/
│   ├── static/
│   ├── package.json
│   └── svelte.config.js
├── web/
│   └── dist/
├── datasets/
│   └── starter/
├── scripts/
├── tests/
├── go.mod
├── Makefile
└── README.md
```

---

## 21. Logging

The Go application should write structured logs containing:

* Startup events.
* Selected port.
* Database migrations.
* Import progress.
* Import errors.
* Score recalculations.
* API errors.
* Shutdown events.

Logs must not include:

* Session tokens.
* Complete imported records containing sensitive data.
* Database credentials.
* Unsafe raw request bodies.

Example log:

```text
time=2026-07-13T10:15:00-06:00
level=INFO
event=server_started
address=127.0.0.1:8787
version=0.1.0
```

---

## 22. Error handling

### Port unavailable

Try the next available port and open that address.

### Database unavailable

Display a local recovery page with:

* Error description.
* Database path.
* Backup option.
* Retry button.
* Reset option requiring confirmation.

### Import failure

Keep the existing valid dataset active. A failed import must not partially replace production data.

### Browser-opening failure

Continue running and print the local URL.

### Frontend asset failure

Return a readable server error page rather than an empty response.

### Unexpected panic

Recover at the HTTP boundary, log the error, and return a generic internal-error response.

---

## 23. Performance targets

For 50 states and up to 100 metrics:

* Application ready within 3 seconds on a typical computer.
* Ranking response under 200 milliseconds after warm-up.
* State comparison response under 200 milliseconds.
* Weight-preview update under 100 milliseconds in the browser.
* Full score recalculation under 2 seconds.
* CSV import progress visible for operations longer than 500 milliseconds.
* Common ranking queries must not require re-reading raw source files.

Scores should be calculated during imports or profile updates and stored as snapshots when practical.

---

## 24. Accessibility requirements

The interface must support:

* Keyboard-only navigation.
* Visible focus indicators.
* Semantic headings.
* Accessible table labels.
* Text alternatives for charts.
* Sufficient contrast.
* Screen-reader descriptions of scores.
* Non-color indicators for better and worse values.
* Reduced-motion preference.
* Responsive layouts.

Charts must supplement tables, not replace them.

---

## 25. Testing strategy

### Go tests

* Metric normalization.
* Category scoring.
* Overall scoring.
* Missing-data behavior.
* Import validation.
* Database migrations.
* API handlers.
* Port-selection behavior.
* Graceful shutdown.
* Host-header validation.

### Frontend tests

* Ranking sorting.
* State search.
* Comparison selection.
* Weight validation.
* Error states.
* Loading states.
* Accessibility checks.

### Integration tests

* Start the compiled local server.
* Verify the frontend loads.
* Verify API communication.
* Import a sample dataset.
* Recalculate scores.
* Export rankings.
* Shut down cleanly.

### Scoring fixture

Create a small test dataset with three fictional states and manually calculated expected results. This makes scoring regressions easy to detect.

---

## 26. MVP scope

The first usable version must include:

* Local Go HTTP server.
* Automatic browser launch.
* Embedded Svelte production build.
* SQLite database.
* All 50 states.
* Five categories.
* Ten total metrics.
* One common comparison year where possible.
* Rankings page.
* State profile page.
* Two-state comparison.
* Adjustable category weights.
* Balanced default profile.
* Methodology page.
* Source attribution.
* CSV export.
* Manual CSV import.
* Graceful shutdown.

The MVP should not attempt to support every available government dataset.

---

## 27. Development phases

### Phase 1: Application shell

* Create Go server.
* Serve a placeholder page.
* Select localhost port.
* Open the browser.
* Implement graceful shutdown.
* Add SQLite and migrations.

### Phase 2: Frontend integration

* Create SvelteKit frontend.
* Configure static build.
* Embed frontend assets in Go.
* Add SPA fallback routing.
* Implement frontend API client.

### Phase 3: State data model

* Add states.
* Add categories.
* Add metrics.
* Add metric values.
* Add sources and imports.

### Phase 4: Scoring engine

* Add normalization.
* Add category scoring.
* Add overall scoring.
* Add missing-data handling.
* Add calculation tests.

### Phase 5: Core interface

* Rankings.
* State profiles.
* Comparison.
* Custom weights.
* Methodology.

### Phase 6: Dataset tools

* CSV import.
* Validation reports.
* Import history.
* Score recalculation.
* Source management.

### Phase 7: Distribution

* Windows build.
* macOS build.
* Linux build.
* Application-data directories.
* Backup and restore.
* Release documentation.

---

## 28. Acceptance criteria

The MVP is complete when:

1. Running one executable starts the local application.
2. The application binds only to localhost.
3. The default browser opens automatically.
4. The ranking page displays all 50 states.
5. Every overall score can be traced to its category and metric values.
6. Users can compare at least two states.
7. Users can change weights and see the ranking change.
8. Missing values are clearly identified.
9. Every metric displays its source and year.
10. The application works offline after data is imported.
11. Restarting the application preserves profiles and settings.
12. A user can export rankings to CSV.
13. A user can import a supported CSV dataset.
14. Failed imports do not corrupt active data.
15. The application shuts down gracefully.
16. Automated tests verify the scoring formulas.

---

## 29. Final product principle

StateScore must answer two separate questions:

```text
What do the public statistics say about each
```

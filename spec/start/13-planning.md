
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

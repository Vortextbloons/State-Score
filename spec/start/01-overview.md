# StateScore â€” Product and Technical Specification

## 1. Product summary

**StateScore** is a locally installed web application for comparing U.S. states across categories such as affordability, economy, education, health, safety, and quality of life.

When the user starts the application:

1. A Go process starts a local HTTP server.
2. The server loads the local SQLite database.
3. The application selects an available localhost port.
4. The userâ€™s default browser opens automatically.
5. The browser displays the Svelte interface.
6. All application data remains on the userâ€™s computer unless the user explicitly downloads updated public datasets.

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

The score must be transparent. The application must never present a stateâ€™s overall score without showing how it was calculated.

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

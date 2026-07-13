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
* Do not serve arbitrary files from the userâ€™s computer.

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

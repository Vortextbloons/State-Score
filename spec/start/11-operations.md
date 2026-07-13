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


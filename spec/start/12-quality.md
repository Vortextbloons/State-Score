
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


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
Data coverage: 2024â€“2026
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


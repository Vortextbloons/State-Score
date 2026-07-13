
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
* Bachelorâ€™s degree attainment.
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
score = ((value - minimum) / (maximum - minimum)) Ã— 100
```

For a lower-is-better metric:

```text
score = ((maximum - value) / (maximum - minimum)) Ã— 100
```

The implementation should support multiple normalization methods:

* Min-max.
* Percentile rank.
* Z-score converted to a 0â€“100 scale.
* Fixed policy thresholds.

Percentile scoring is recommended as the default because basic min-max scoring can be distorted by extreme values.

### 11.2 Category score

```text
category score =
sum(metric score Ã— metric weight)
/
sum(included metric weights)
```

### 11.3 Overall score

```text
overall score =
economy score Ã— economy weight
+ education score Ã— education weight
+ health score Ã— health weight
+ safety score Ã— safety weight
+ affordability score Ã— affordability weight
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

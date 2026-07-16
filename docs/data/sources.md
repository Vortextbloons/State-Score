# Data Sources

This is the canonical directory of the public datasets behind StateScore's active metrics. Source metadata is also stored in `data_sources`, and each observation links through `metric_values.import_id` to its import record.

State reference data also includes the Census Bureau's Vintage 2025 resident population estimate for July 1, 2025. Population is descriptive and does not participate in scoring. Values come from [`NST-EST2025-ALLDATA.csv`](https://www2.census.gov/programs-surveys/popest/datasets/2020-2025/state/totals/NST-EST2025-ALLDATA.csv); each state records the estimate, year, and source ID.

## Source Catalog

| Category | Metric | Year | Publisher and dataset | Where to find it |
|---|---|---:|---|---|
| Economy | `unemployment-rate` | 2024 | U.S. Census Bureau, ACS state indicators, table B23025 | [ACS 2024 5-Year API](https://api.census.gov/data/2024/acs/acs5) |
| Economy | `median-household-income` | 2024 | U.S. Census Bureau, ACS table B19013 | [ACS 2024 5-Year API](https://api.census.gov/data/2024/acs/acs5) |
| Economy | `annual-employment-growth` | 2024 | Bureau of Labor Statistics, Current Employment Statistics State and Metro Area | [BLS Public Data API](https://api.bls.gov/publicAPI/v2/timeseries/data/) and [CES series structure](https://www.bls.gov/sae/additional-resources/state-and-area-ces-series-code-structure-under-naics.htm) |
| Economy | `labor-force-participation-rate` | 2024 | U.S. Census Bureau, ACS 1-Year Subject Table S2301 | [ACS Subject API](https://api.census.gov/data/2024/acs/acs1/subject) |
| Education | `high-school-graduation-rate` | 2024 | U.S. Census Bureau, ACS table B15003 | [ACS 2024 5-Year API](https://api.census.gov/data/2024/acs/acs5) |
| Education | `bachelors-degree-attainment` | 2024 | U.S. Census Bureau, ACS table B15003 | [ACS 2024 5-Year API](https://api.census.gov/data/2024/acs/acs5) |
| Education | `young-adult-college-enrollment` | 2024 | U.S. Census Bureau, ACS 1-Year Subject Table S1401 | [ACS Subject API](https://api.census.gov/data/2024/acs/acs1/subject) |
| Education | `naep-achievement-composite` | 2024 | NCES, National Assessment of Educational Progress reading and mathematics | [Nation's Report Card 2024 reports](https://www.nationsreportcard.gov/reports/) |
| Health | `life-expectancy` | 2022 | CDC/NCHS, U.S. State Life Tables | [National Vital Statistics Report 74-12](https://www.cdc.gov/nchs/data/nvsr/nvsr74/nvsr74-12.pdf) |
| Health | `adult-obesity-prevalence` | 2024 | CDC, BRFSS Nutrition, Physical Activity and Obesity dataset `hn4x-zwk7` | [Data.CDC.gov dataset](https://data.cdc.gov/d/hn4x-zwk7) and [Socrata API](https://data.cdc.gov/resource/hn4x-zwk7.json) |
| Health | `uninsured-rate` | 2024 | U.S. Census Bureau, ACS 1-Year Subject Table S2701 | [ACS Subject API](https://api.census.gov/data/2024/acs/acs1/subject) |
| Safety | `violent-crime-rate` | 2024 | FBI, Uniform Crime Reporting / Crime Data Explorer | [Crime Data Explorer](https://cde.ucr.cjis.gov/LATEST/webapp/#/pages/explorer/crime/crime-trend) |
| Safety | `traffic-fatalities` | 2024 | NHTSA FARS, tabulated by the Insurance Institute for Highway Safety | [IIHS state-by-state fatality statistics](https://www.iihs.org/research-areas/fatality-statistics/detail/state-by-state) |
| Safety | `property-crime-rate` | 2024 | FBI, Uniform Crime Reporting / Crime Data Explorer summarized state data | [CDE summarized-state API base](https://cde.ucr.cjis.gov/LATEST/summarized/state) |
| Safety | `age-adjusted-homicide-death-rate` | 2024 | CDC/NCHS, National Vital Statistics System via CDC WONDER | [CDC homicide mortality table](https://www.cdc.gov/nchs/state-stats/deaths/homicide.html) |
| Affordability | `cost-of-living-index` | 2024 | Bureau of Economic Analysis, Regional Price Parities by State | [BEA SARPP download](https://apps.bea.gov/regional/zip/SARPP.zip) |
| Affordability | `renter-housing-cost-burden` | 2024 | U.S. Census Bureau, ACS 1-Year Data Profile DP04 | [ACS Profile API](https://api.census.gov/data/2024/acs/acs1/profile) |
| Affordability | `owner-housing-cost-burden` | 2024 | U.S. Census Bureau, ACS 1-Year Detailed Table B25091 | [ACS Detailed API](https://api.census.gov/data/2024/acs/acs1) |

All listed sources are U.S. government public data except the IIHS presentation of NHTSA FARS data.

## Retrieval and Calculations

### ACS economy and attainment metrics

The four original Economy and Education observations come from the 2024 ACS 5-Year API. Their source tables are B23025, B19013, and B15003. StateScore calculates unemployment as unemployed residents divided by the civilian labor force; educational attainment uses the population age 25 and older.

API base:

```text
https://api.census.gov/data/2024/acs/acs5
```

### Annual employment growth

Use seasonally adjusted statewide total nonfarm CES series. The series ID format is `SMS{state_fips}000000000000001`. Retrieve all twelve monthly observations for 2023 and 2024, average each calendar year, then calculate:

```text
((2024 annual average / 2023 annual average) - 1) * 100
```

API base:

```text
POST https://api.bls.gov/publicAPI/v2/timeseries/data/
```

### Labor-force participation

Use `S2301_C02_001E`, the labor-force participation rate for the population age 16 and over, from the 2024 ACS 1-Year Subject API.

### NAEP achievement composite

For each state, take the arithmetic mean of the 2024 average scale scores for public-school students in grade 4 mathematics, grade 8 mathematics, grade 4 reading, and grade 8 reading. All four assessments use the NAEP 0–500 scale. StateScore stores the two-decimal composite before percentile normalization.

### Uninsured rate

Use `S2701_C05_001E`, the percent uninsured among the civilian noninstitutionalized population, from the 2024 ACS 1-Year Subject API.

### Young-adult college enrollment

Variables:

- `S1401_C01_029E`: total population ages 18-24
- `S1401_C01_030E`: ages 18-24 enrolled in college or graduate school

```text
S1401_C01_030E / S1401_C01_029E * 100
```

Request:

```text
https://api.census.gov/data/2024/acs/acs1/subject
  ?get=NAME,S1401_C01_029E,S1401_C01_030E
  &for=state:*
  &key=YOUR_CENSUS_API_KEY
```

### Adult obesity prevalence

Query dataset `hn4x-zwk7` for year 2024, question ID `Q036`, and `stratificationcategory1='Total'`. Use only the overall-population crude estimate. Tennessee's 2024 overall value is suppressed by the source, so StateScore stores 49 observations and treats Tennessee as missing.

### Property-crime rate

For each postal abbreviation, request all months in 2024:

```text
https://cde.ucr.cjis.gov/LATEST/summarized/state/{STATE}/property-crime
  ?from=01-2024
  &to=12-2024
```

StateScore sums monthly offenses and divides by the FBI state population:

```text
annual offenses / state population * 100000
```

The import stores the minimum monthly reporting coverage, average participating population, and CDE revision date. Observations below 90% minimum monthly population coverage remain visible but are excluded from scoring. The FBI response does not currently supply a participating-agency count, so that nullable provenance field is retained for future source revisions.

When the latest year fails the gate, the as-of scorer may use the newest older observation that independently clears the same 90% threshold. The bundled fallbacks are Louisiana 2023, Florida 2020, and South Dakota 2020. Indiana, Mississippi, and Wyoming have no qualifying fallback in the years reviewed and remain incomplete.

### Renter housing-cost burden

Variables:

- `DP04_0141PE`: renter households spending 30.0%-34.9% of income on gross rent
- `DP04_0142PE`: renter households spending 35.0% or more

```text
DP04_0141PE + DP04_0142PE
```

Request:

```text
https://api.census.gov/data/2024/acs/acs1/profile
  ?get=NAME,DP04_0141PE,DP04_0142PE
  &for=state:*
  &key=YOUR_CENSUS_API_KEY
```

## Credentials

The Census Data API requires a free API key. Store it only in the Git-ignored root `.env` file:

```text
CENSUS_API_KEY=your_key_here
```

Never commit the key or place its value in documentation, migrations, logs, or request examples.

### Age-adjusted homicide death rate

Use the final 2024 NVSS age-adjusted rate published by CDC's Stats of the States table, sourced from CDC WONDER. Homicide uses underlying-cause codes `*U01–*U02`, `X85–Y09`, and `Y87.1`; rates are age-adjusted to the 2000 U.S. standard population.

### Owner housing-cost burden

From ACS detailed table B25091, the denominator is the sum of computable selected-monthly-owner-cost categories `003–010` and `014–021`. The numerator is the six categories at or above 30 percent: `008–010` and `019–021`.

```text
(B25091_008E + B25091_009E + B25091_010E +
 B25091_019E + B25091_020E + B25091_021E)
/
(sum of B25091_003E:B25091_010E and B25091_014E:B25091_021E) * 100
```

## Bundled Snapshots

The reproducible source metadata, import checksums, and bundled values live in these migrations:

- `000006_seed_safety_data.sql`
- `000007_seed_health_affordability_data.sql`
- `000009_seed_economy_education_data.sql`
- `000010_add_priority_metrics.sql`
- `000011_add_foundational_metrics.sql`
- `000012_add_state_population.sql`

When refreshing a source, add a new migration or use the managed import workflow; do not rewrite an already-applied migration.

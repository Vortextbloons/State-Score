-- Fix construct labels and stop empty metrics from diluting "balanced" rankings.
-- Metrics without seeded coverage stay defined but inactive until data is imported.

UPDATE metrics
SET name = 'High school attainment',
    description = 'Percentage of adults age 25+ with a high school diploma or equivalent (attainment, not cohort graduation)',
    updated_at = datetime('now')
WHERE slug = 'high-school-graduation-rate';

UPDATE metrics
SET name = 'Regional price parity',
    description = 'BEA all-items regional price parity index (U.S. average = 100). Higher values mean prices above the national average.',
    unit = 'Index (US=100)',
    updated_at = datetime('now')
WHERE slug = 'cost-of-living-index';

UPDATE metrics
SET active = 0,
    updated_at = datetime('now')
WHERE slug IN (
    'unemployment-rate',
    'median-household-income',
    'high-school-graduation-rate',
    'bachelors-degree-attainment',
    'uninsured-rate',
    'median-rent'
);

UPDATE scoring_profiles
SET description = 'Equal weight across categories that have active metrics. Inactive metrics are excluded until data is imported.',
    updated_at = datetime('now')
WHERE name = 'Balanced';

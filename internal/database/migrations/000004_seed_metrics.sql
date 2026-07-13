-- Seed initial metrics (2 per category for MVP).

-- Economy metrics
INSERT INTO metrics (category_id, slug, name, description, unit, higher_is_better, normalization_method, default_weight)
SELECT id, 'unemployment-rate', 'Unemployment rate', 'Percentage of labor force that is unemployed', 'Percent', 0, 'percentile', 0.50
FROM categories WHERE slug = 'economy';

INSERT INTO metrics (category_id, slug, name, description, unit, higher_is_better, normalization_method, default_weight)
SELECT id, 'median-household-income', 'Median household income', 'Annual median household income in USD', 'Dollars', 1, 'percentile', 0.50
FROM categories WHERE slug = 'economy';

-- Education metrics
INSERT INTO metrics (category_id, slug, name, description, unit, higher_is_better, normalization_method, default_weight)
SELECT id, 'high-school-graduation-rate', 'High school graduation rate', 'Percentage of adults with a high school diploma or equivalent', 'Percent', 1, 'percentile', 0.50
FROM categories WHERE slug = 'education';

INSERT INTO metrics (category_id, slug, name, description, unit, higher_is_better, normalization_method, default_weight)
SELECT id, 'bachelors-degree-attainment', 'Bachelor''s degree attainment', 'Percentage of adults age 25+ with a bachelor''s degree or higher', 'Percent', 1, 'percentile', 0.50
FROM categories WHERE slug = 'education';

-- Health metrics
INSERT INTO metrics (category_id, slug, name, description, unit, higher_is_better, normalization_method, default_weight)
SELECT id, 'life-expectancy', 'Life expectancy', 'Average life expectancy at birth in years', 'Years', 1, 'percentile', 0.50
FROM categories WHERE slug = 'health';

INSERT INTO metrics (category_id, slug, name, description, unit, higher_is_better, normalization_method, default_weight)
SELECT id, 'uninsured-rate', 'Uninsured rate', 'Percentage of population without health insurance', 'Percent', 0, 'percentile', 0.50
FROM categories WHERE slug = 'health';

-- Safety metrics
INSERT INTO metrics (category_id, slug, name, description, unit, higher_is_better, normalization_method, default_weight)
SELECT id, 'violent-crime-rate', 'Violent crime rate', 'Violent crimes per 100,000 population', 'Per 100k', 0, 'percentile', 0.50
FROM categories WHERE slug = 'safety';

INSERT INTO metrics (category_id, slug, name, description, unit, higher_is_better, normalization_method, default_weight)
SELECT id, 'traffic-fatalities', 'Traffic fatalities', 'Traffic fatalities per 100,000 population', 'Per 100k', 0, 'percentile', 0.50
FROM categories WHERE slug = 'safety';

-- Affordability metrics
INSERT INTO metrics (category_id, slug, name, description, unit, higher_is_better, normalization_method, default_weight)
SELECT id, 'median-rent', 'Median rent', 'Median monthly gross rent in USD', 'Dollars', 0, 'percentile', 0.50
FROM categories WHERE slug = 'affordability';

INSERT INTO metrics (category_id, slug, name, description, unit, higher_is_better, normalization_method, default_weight)
SELECT id, 'cost-of-living-index', 'Cost of living index', 'Regional cost of living index (US average = 100)', 'Index', 0, 'percentile', 0.50
FROM categories WHERE slug = 'affordability';

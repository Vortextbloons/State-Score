-- Add the first expansion metric for each scoring category.
-- All five metrics include a bundled 2024 snapshot from their official APIs.

CREATE TABLE metric_value_quality (
    metric_value_id       INTEGER PRIMARY KEY REFERENCES metric_values(id) ON DELETE CASCADE,
    reporting_coverage    REAL,
    participating_agencies INTEGER,
    population_covered    INTEGER,
    data_revision         TEXT,
    scoring_eligible      INTEGER NOT NULL DEFAULT 1,
    exclusion_reason      TEXT
);

INSERT INTO data_sources(name,publisher,source_url,license,format,description) VALUES
('CES State and Metro Area 2024','U.S. Bureau of Labor Statistics','https://api.bls.gov/publicAPI/v2/timeseries/data/','U.S. government public data','json','Seasonally adjusted statewide total nonfarm employment; annual growth compares 2024 and 2023 monthly annual averages.'),
('ACS 2024 Subject Table S1401','U.S. Census Bureau','https://api.census.gov/data/2024/acs/acs1/subject','U.S. government public data','json','Population age 18-24 enrolled in college or graduate school divided by the total population age 18-24.'),
('CDC BRFSS Nutrition, Physical Activity and Obesity 2024','Centers for Disease Control and Prevention','https://data.cdc.gov/resource/hn4x-zwk7.json','U.S. government public data','json','Overall-population crude adult obesity prevalence; stratified estimates are excluded.'),
('FBI Crime Data Explorer 2024 Property Crime','Federal Bureau of Investigation','https://cde.ucr.cjis.gov/LATEST/summarized/state','U.S. government public data','json','Annual burglary, larceny-theft, and motor-vehicle-theft offenses per 100,000; states below 90% minimum monthly population coverage are excluded from scoring.'),
('ACS 2024 Data Profile DP04','U.S. Census Bureau','https://api.census.gov/data/2024/acs/acs1/profile','U.S. government public data','json','Renter households spending at least 30 percent of household income on gross rent.');

WITH v(category_slug,slug,name,description,unit,higher,source_name,active) AS (VALUES
 ('economy','annual-employment-growth','Annual employment growth','Percentage change in annual-average seasonally adjusted total nonfarm employment from the previous year.','Percent',1,'CES State and Metro Area 2024',1),
 ('education','young-adult-college-enrollment','Young-adult college enrollment rate','Percentage of residents ages 18-24 enrolled in college or graduate school.','Percent',1,'ACS 2024 Subject Table S1401',1),
 ('health','adult-obesity-prevalence','Adult obesity prevalence','Percentage of adults classified as having obesity from self-reported height and weight.','Percent',0,'CDC BRFSS Nutrition, Physical Activity and Obesity 2024',1),
 ('safety','property-crime-rate','Property-crime rate','Reported burglary, larceny-theft, and motor-vehicle-theft offenses per 100,000 residents.','Per 100k',0,'FBI Crime Data Explorer 2024 Property Crime',1),
 ('affordability','renter-housing-cost-burden','Renter housing-cost burden','Percentage of renter households spending at least 30 percent of household income on gross rent.','Percent',0,'ACS 2024 Data Profile DP04',1)
)
INSERT INTO metrics(category_id,slug,name,description,unit,higher_is_better,normalization_method,default_weight,source_id,active)
SELECT c.id,v.slug,v.name,v.description,v.unit,v.higher,'percentile',1.0,ds.id,v.active
FROM v
JOIN categories c ON c.slug=v.category_slug
JOIN data_sources ds ON ds.name=v.source_name;

-- Existing metrics and the new metric share equal defaults within each category.
UPDATE metrics SET default_weight=1.0/(
  SELECT count(*) FROM metrics sibling WHERE sibling.category_id=metrics.category_id AND sibling.active=1
) WHERE active=1;

INSERT OR REPLACE INTO profile_metric_weights(profile_id,metric_id,weight)
SELECT p.id,m.id,m.default_weight FROM scoring_profiles p JOIN metrics m ON m.active=1;

INSERT INTO imports(source_id,status,started_at,completed_at,records_read,records_inserted,records_rejected,checksum)
SELECT id,'completed',datetime('now'),datetime('now'),50,50,0,'bundled-priority-metrics-2024-v1' FROM data_sources WHERE name='CES State and Metro Area 2024';
INSERT INTO imports(source_id,status,started_at,completed_at,records_read,records_inserted,records_rejected,checksum)
SELECT id,'completed_with_errors',datetime('now'),datetime('now'),50,49,1,'bundled-priority-obesity-2024-v1' FROM data_sources WHERE name='CDC BRFSS Nutrition, Physical Activity and Obesity 2024';
INSERT INTO imports(source_id,status,started_at,completed_at,records_read,records_inserted,records_rejected,checksum)
SELECT id,'completed',datetime('now'),datetime('now'),50,50,0,'bundled-priority-property-crime-2024-v1' FROM data_sources WHERE name='FBI Crime Data Explorer 2024 Property Crime';
INSERT INTO imports(source_id,status,started_at,completed_at,records_read,records_inserted,records_rejected,checksum)
SELECT id,'completed',datetime('now'),datetime('now'),50,50,0,'bundled-priority-college-enrollment-2024-v1' FROM data_sources WHERE name='ACS 2024 Subject Table S1401';
INSERT INTO imports(source_id,status,started_at,completed_at,records_read,records_inserted,records_rejected,checksum)
SELECT id,'completed',datetime('now'),datetime('now'),50,50,0,'bundled-priority-renter-burden-2024-v1' FROM data_sources WHERE name='ACS 2024 Data Profile DP04';

WITH o(state_name,college_enrollment,renter_burden) AS (VALUES
('Alabama',43.7307,48.2),('Alaska',18.7293,41.4),('Arizona',36.8557,52.4),('Arkansas',31.7274,45.5),('California',47.0084,55.8),
('Colorado',39.2827,52.5),('Connecticut',47.9917,54.2),('Delaware',44.59,50.3),('Florida',40.8444,62.1),('Georgia',37.7365,53.1),
('Hawaii',34.3776,55),('Idaho',37.8976,48.4),('Illinois',41.5383,48.7),('Indiana',35.9942,48.5),('Iowa',41.188,43.3),
('Kansas',39.343,43.4),('Kentucky',36.236,45.2),('Louisiana',35.4229,54.2),('Maine',47.2727,50.2),('Maryland',41.2143,51.6),
('Massachusetts',50.5905,52.4),('Michigan',40.2664,51.2),('Minnesota',37.5445,47.7),('Mississippi',37.589,50.3),('Missouri',35.511,45.2),
('Montana',30.0418,42.7),('Nebraska',41.6984,46.5),('Nevada',28.904,57.6),('New Hampshire',42.9739,49.4),('New Jersey',44.5946,52.3),
('New Mexico',33.1668,48.4),('New York',47.5127,51.5),('North Carolina',39.4075,51.6),('North Dakota',46.2044,36.2),('Ohio',37.4814,47.7),
('Oklahoma',33.7923,46.1),('Oregon',36.8753,52.3),('Pennsylvania',44.029,49.4),('Rhode Island',50.6079,51.1),('South Carolina',37.19,52.6),
('South Dakota',30.8881,41.9),('Tennessee',33.5964,50.1),('Texas',35.869,52.6),('Utah',43.1326,48.2),('Vermont',45.2027,48.8),
('Virginia',42.7489,49.3),('Washington',33.2875,50.7),('West Virginia',34.5711,47),('Wisconsin',38.7202,44.9),('Wyoming',39.0069,42.9)
), expanded(state_name,slug,value) AS (
 SELECT state_name,'young-adult-college-enrollment',college_enrollment FROM o
 UNION ALL SELECT state_name,'renter-housing-cost-burden',renter_burden FROM o
)
INSERT INTO metric_values(state_id,metric_id,year,value,source_record_id,import_id)
SELECT s.id,m.id,2024,e.value,e.state_name,i.id FROM expanded e
JOIN states s ON s.name=e.state_name JOIN metrics m ON m.slug=e.slug
JOIN imports i ON i.checksum=CASE e.slug WHEN 'young-adult-college-enrollment' THEN 'bundled-priority-college-enrollment-2024-v1' ELSE 'bundled-priority-renter-burden-2024-v1' END;

WITH o(code,employment,obesity,crime,coverage,population_covered) AS (VALUES
('AK',2.009,34,1741.04,99.67,737670),('AL',1.2855,38.9,1624.62,98.1,5104327),('AR',1.377,38.9,1940.35,98.12,3041532),('AZ',1.9511,33.3,1788.11,99.29,7537200),('CA',0.8792,29.1,1986.66,99.46,39305752),
('CO',1.0959,25,2641.95,98.78,5911687),('CT',0.8828,32,1397.33,100,3675069),('DE',1.578,36.6,1773.62,99.89,1051804),('FL',1.7104,29.6,1035.42,78.32,18788884),('GA',1.1648,35.4,1582.77,91.82,10355890),
('HI',0.962,27,2057.95,100,1446146),('IA',0.3499,36.6,1286.63,95.61,3133294),('ID',1.8726,32.7,755.24,99.07,1994351),('IL',0.4089,34.2,1665.69,93.34,11991926),('IN',0.6329,38.4,1313.9,86,6040464),
('KS',0.9911,37.6,2089.98,95.61,2855099),('KY',0.9934,37.2,1400.24,100,4588372),('LA',1.331,39.2,1781.05,80.3,3771756),('MA',0.31,27,1122.94,99.63,7111694),('MD',2.2002,32.7,2076.2,100,6263220),
('ME',1.3746,33.2,1149.24,99.28,1396149),('MI',0.5825,36.1,1399.69,97.41,9934280),('MN',1.0426,32.3,1625.14,99.81,5786491),('MO',0.6792,34.6,1973.21,97.87,6138924),('MS',0.5488,40.4,952.89,64.26,1915720),
('MT',1.1482,31,1648.12,99.92,1136348),('NC',1.4812,34.5,1948.37,96.38,10696964),('ND',1.8571,36.8,1707.07,99.85,795372),('NE',0.9288,37.6,1632.14,96.56,1952076),('NH',0.7359,31.1,928.08,99.86,1408213),
('NJ',0.901,27.7,1398.58,94.88,9271582),('NM',1.5196,34.5,2707.33,92.04,1994276),('NV',1.809,34.2,2226.68,99.41,3252048),('NY',1.7076,29.5,1667.1,91.42,19275838),('OH',0.7104,36.9,1551.82,93.81,11200167),
('OK',1.3478,36.8,2009.5,99.76,4089788),('OR',0.2425,33.5,2413.46,98.04,4204673),('PA',0.8766,34.2,1457.87,98.88,13004606),('RI',1.354,31.1,1040.18,99.87,1110883),('SC',1.946,34.6,1988.56,99.49,5464631),
('SD',1.3136,37,1535.47,89.19,828963),('TN',1.1512,NULL,2058.37,99.39,7200395),('TX',1.6617,35.6,2059.47,99.31,31083196),('UT',1.3508,31,1436.03,91.53,3357699),('VA',1.4723,32.3,1591,99.63,8803162),
('VT',0.6988,29,1667.25,100,648493),('WA',0.7156,31.5,2505.95,99.37,7914020),('WI',0.7314,37.4,1164.42,95.42,5821643),('WV',1.3094,41.4,1089.79,91.1,1629954),('WY',1.0902,32.5,1188.53,88.89,528544)
), expanded(code,slug,value) AS (
 SELECT code,'annual-employment-growth',employment FROM o
 UNION ALL SELECT code,'adult-obesity-prevalence',obesity FROM o WHERE obesity IS NOT NULL
 UNION ALL SELECT code,'property-crime-rate',crime FROM o
)
INSERT INTO metric_values(state_id,metric_id,year,value,source_record_id,import_id)
SELECT s.id,m.id,2024,e.value,e.code,i.id FROM expanded e
JOIN states s ON s.code=e.code JOIN metrics m ON m.slug=e.slug
JOIN imports i ON i.checksum=CASE e.slug WHEN 'annual-employment-growth' THEN 'bundled-priority-metrics-2024-v1' WHEN 'adult-obesity-prevalence' THEN 'bundled-priority-obesity-2024-v1' ELSE 'bundled-priority-property-crime-2024-v1' END;

WITH q(code,coverage,population_covered) AS (VALUES
('AK',99.67,737670),('AL',98.1,5104327),('AR',98.12,3041532),('AZ',99.29,7537200),('CA',99.46,39305752),('CO',98.78,5911687),('CT',100,3675069),('DE',99.89,1051804),('FL',78.32,18788884),('GA',91.82,10355890),('HI',100,1446146),('IA',95.61,3133294),('ID',99.07,1994351),('IL',93.34,11991926),('IN',86,6040464),('KS',95.61,2855099),('KY',100,4588372),('LA',80.3,3771756),('MA',99.63,7111694),('MD',100,6263220),('ME',99.28,1396149),('MI',97.41,9934280),('MN',99.81,5786491),('MO',97.87,6138924),('MS',64.26,1915720),('MT',99.92,1136348),('NC',96.38,10696964),('ND',99.85,795372),('NE',96.56,1952076),('NH',99.86,1408213),('NJ',94.88,9271582),('NM',92.04,1994276),('NV',99.41,3252048),('NY',91.42,19275838),('OH',93.81,11200167),('OK',99.76,4089788),('OR',98.04,4204673),('PA',98.88,13004606),('RI',99.87,1110883),('SC',99.49,5464631),('SD',89.19,828963),('TN',99.39,7200395),('TX',99.31,31083196),('UT',91.53,3357699),('VA',99.63,8803162),('VT',100,648493),('WA',99.37,7914020),('WI',95.42,5821643),('WV',91.1,1629954),('WY',88.89,528544)
)
INSERT INTO metric_value_quality(metric_value_id,reporting_coverage,population_covered,data_revision,scoring_eligible,exclusion_reason)
SELECT mv.id,q.coverage,q.population_covered,'2026-07-15',q.coverage>=90,CASE WHEN q.coverage<90 THEN 'Minimum monthly FBI population coverage below 90%' END
FROM q JOIN states s ON s.code=q.code JOIN metrics m ON m.slug='property-crime-rate'
JOIN metric_values mv ON mv.state_id=s.id AND mv.metric_id=m.id AND mv.year=2024;

-- Coverage-qualified fallbacks used by the as-of scorer when a 2024 observation
-- is suppressed or fails the FBI quality gate. No fallback is bundled unless it
-- independently satisfies the same source/methodology rules.
INSERT INTO imports(source_id,status,started_at,completed_at,records_read,records_inserted,records_rejected,checksum)
SELECT id,'completed',datetime('now'),datetime('now'),3,3,0,'bundled-priority-property-crime-fallbacks-v1' FROM data_sources WHERE name='FBI Crime Data Explorer 2024 Property Crime';
INSERT INTO imports(source_id,status,started_at,completed_at,records_read,records_inserted,records_rejected,checksum)
SELECT id,'completed',datetime('now'),datetime('now'),1,1,0,'bundled-priority-obesity-fallback-v1' FROM data_sources WHERE name='CDC BRFSS Nutrition, Physical Activity and Obesity 2024';

WITH o(code,year,value,coverage,population_covered) AS (VALUES
 ('FL',2020,1768.08,99.57,21665734),
 ('LA',2023,2437.83,91.48,4256322),
 ('SD',2020,1821.63,94.76,847803)
)
INSERT INTO metric_values(state_id,metric_id,year,value,source_record_id,import_id)
SELECT s.id,m.id,o.year,o.value,o.code,i.id FROM o
JOIN states s ON s.code=o.code JOIN metrics m ON m.slug='property-crime-rate'
JOIN imports i ON i.checksum='bundled-priority-property-crime-fallbacks-v1';

WITH q(code,year,coverage,population_covered) AS (VALUES
 ('FL',2020,99.57,21665734),('LA',2023,91.48,4256322),('SD',2020,94.76,847803)
)
INSERT INTO metric_value_quality(metric_value_id,reporting_coverage,population_covered,data_revision,scoring_eligible)
SELECT mv.id,q.coverage,q.population_covered,'2026-07-15',1 FROM q
JOIN states s ON s.code=q.code JOIN metrics m ON m.slug='property-crime-rate'
JOIN metric_values mv ON mv.state_id=s.id AND mv.metric_id=m.id AND mv.year=q.year;

INSERT INTO metric_values(state_id,metric_id,year,value,source_record_id,import_id)
SELECT s.id,m.id,2023,37.6,'TN',i.id FROM states s,metrics m,imports i
WHERE s.code='TN' AND m.slug='adult-obesity-prevalence' AND i.checksum='bundled-priority-obesity-fallback-v1';

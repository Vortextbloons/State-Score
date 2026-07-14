-- Restore Economy and Education with complete 2024 ACS coverage for all 50 states.
-- Source tables: B23025 (employment), B19013 (income), and B15003 (attainment).

INSERT INTO data_sources(name,publisher,source_url,license,format,description) VALUES
('American Community Survey 2024 State Indicators','U.S. Census Bureau','https://www.census.gov/data/developers/data-sets/acs-5year/2024.html','U.S. government public data','json','2024 ACS state estimates. Unemployment is unemployed / civilian labor force; attainment percentages use population age 25+.' );

INSERT INTO imports(source_id,status,started_at,completed_at,records_read,records_inserted,records_rejected,checksum)
SELECT id,'completed',datetime('now'),datetime('now'),200,200,0,'bundled-acs-economy-education-2024-v1'
FROM data_sources WHERE name='American Community Survey 2024 State Indicators';

UPDATE metrics
SET source_id=(SELECT id FROM data_sources WHERE name='American Community Survey 2024 State Indicators'),
    active=1,
    updated_at=datetime('now')
WHERE slug IN ('unemployment-rate','median-household-income','high-school-graduation-rate','bachelors-degree-attainment');

WITH observations(state_name,unemployment,income,high_school,bachelors) AS (VALUES
('Alabama',4.23,66659,89.55,29.85),('Alaska',5.96,95665,93.41,32.75),
('Arizona',4.56,81486,89.74,34.73),('Arkansas',3.82,62106,89.29,27.12),
('California',5.93,100149,84.85,38.12),('Colorado',4.28,97113,93.16,47.77),
('Connecticut',4.95,96049,91.81,42.57),('Delaware',5.11,87534,91.82,35.99),
('Florida',4.56,77735,90.36,35.84),('Georgia',4.69,79991,89.78,36.31),
('Hawaii',3.38,100745,93.71,37.76),('Idaho',3.73,81166,92.49,33.03),
('Illinois',5.13,83211,90.72,39.23),('Indiana',4.10,71959,90.58,30.66),
('Iowa',3.18,75501,92.99,32.05),('Kansas',3.67,75514,92.00,36.05),
('Kentucky',4.87,64526,88.98,27.94),('Louisiana',5.15,60986,88.23,27.84),
('Maine',2.72,76442,94.77,37.11),('Maryland',4.35,102905,91.42,44.69),
('Massachusetts',4.07,104828,91.40,48.27),('Michigan',4.50,72389,92.24,33.33),
('Minnesota',3.72,87117,94.20,40.03),('Mississippi',4.60,59127,88.00,27.02),
('Missouri',3.25,71589,92.01,33.47),('Montana',2.96,75340,94.57,36.34),
('Nebraska',3.40,76376,92.65,35.36),('Nevada',6.28,81134,87.70,28.53),
('New Hampshire',2.78,99782,95.05,41.52),('New Jersey',5.41,104294,90.36,44.55),
('New Mexico',5.27,67816,87.87,31.78),('New York',5.32,85820,88.04,41.22),
('North Carolina',3.99,73958,90.85,37.14),('North Dakota',1.91,77871,94.15,33.95),
('Ohio',3.99,72212,92.18,32.35),('Oklahoma',4.68,66148,90.22,29.32),
('Oregon',4.85,85220,92.30,37.84),('Pennsylvania',4.35,77545,92.35,36.37),
('Rhode Island',5.04,83504,90.42,39.00),('South Carolina',4.50,72350,90.95,33.33),
('South Dakota',2.81,76881,94.03,34.18),('Tennessee',4.23,71997,90.57,32.42),
('Texas',4.83,79721,86.69,35.16),('Utah',3.78,96658,93.62,39.12),
('Vermont',2.67,82730,94.71,45.10),('Virginia',3.62,92090,91.71,43.34),
('Washington',4.87,99389,92.27,41.02),('West Virginia',4.63,60798,90.07,24.38),
('Wisconsin',3.05,77488,93.71,34.56),('Wyoming',3.43,75532,94.42,32.28)
), expanded(state_name,metric_slug,value) AS (
SELECT state_name,'unemployment-rate',unemployment FROM observations UNION ALL
SELECT state_name,'median-household-income',income FROM observations UNION ALL
SELECT state_name,'high-school-graduation-rate',high_school FROM observations UNION ALL
SELECT state_name,'bachelors-degree-attainment',bachelors FROM observations
)
INSERT INTO metric_values(state_id,metric_id,year,value,source_record_id,import_id)
SELECT s.id,m.id,2024,e.value,e.state_name,i.id
FROM expanded e
JOIN states s ON s.name=e.state_name
JOIN metrics m ON m.slug=e.metric_slug
JOIN imports i ON i.checksum='bundled-acs-economy-education-2024-v1';

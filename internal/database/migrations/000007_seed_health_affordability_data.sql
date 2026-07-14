-- Latest complete state-level data for the previously empty Life expectancy
-- and Cost of living index metrics, available as of July 2026.

INSERT INTO data_sources(name,publisher,source_url,license,format,description) VALUES
('U.S. State Life Tables 2022','Centers for Disease Control and Prevention, National Center for Health Statistics','https://www.cdc.gov/nchs/data/nvsr/nvsr74/nvsr74-12.pdf','U.S. government public data','pdf','2022 life expectancy at birth for the total population in each state.'),
('Regional Price Parities by State 2024','U.S. Bureau of Economic Analysis','https://apps.bea.gov/regional/zip/SARPP.zip','U.S. government public data','csv','2024 all-items regional price parity index for each state; U.S. average equals 100.');

INSERT INTO imports(source_id,status,started_at,completed_at,records_read,records_inserted,records_rejected,checksum)
SELECT id,'completed',datetime('now'),datetime('now'),50,50,0,'bundled-cdc-life-2022-v1' FROM data_sources WHERE name='U.S. State Life Tables 2022';
INSERT INTO imports(source_id,status,started_at,completed_at,records_read,records_inserted,records_rejected,checksum)
SELECT id,'completed',datetime('now'),datetime('now'),50,50,0,'bundled-bea-rpp-2024-v1' FROM data_sources WHERE name='Regional Price Parities by State 2024';

UPDATE metrics SET source_id=(SELECT id FROM data_sources WHERE name='U.S. State Life Tables 2022') WHERE slug='life-expectancy';
UPDATE metrics SET source_id=(SELECT id FROM data_sources WHERE name='Regional Price Parities by State 2024') WHERE slug='cost-of-living-index';

WITH observations(state_name,value) AS (VALUES
('Alabama',73.8),('Alaska',75.8),('Arizona',76.7),('Arkansas',73.9),('California',79.3),
('Colorado',78.5),('Connecticut',79.4),('Delaware',76.5),('Florida',77.9),('Georgia',75.9),
('Hawaii',80.0),('Idaho',78.4),('Illinois',77.5),('Indiana',75.4),('Iowa',77.9),
('Kansas',76.5),('Kentucky',73.6),('Louisiana',73.8),('Maine',76.6),('Maryland',77.8),
('Massachusetts',79.8),('Michigan',76.8),('Minnesota',79.3),('Mississippi',72.6),('Missouri',75.2),
('Montana',77.3),('Nebraska',78.3),('Nevada',76.4),('New Hampshire',78.7),('New Jersey',79.6),
('New Mexico',74.5),('New York',79.5),('North Carolina',75.9),('North Dakota',77.9),('Ohio',75.6),
('Oklahoma',73.8),('Oregon',77.7),('Pennsylvania',77.3),('Rhode Island',79.2),('South Carolina',75.1),
('South Dakota',77.3),('Tennessee',73.8),('Texas',77.1),('Utah',79.0),('Vermont',78.3),
('Virginia',77.5),('Washington',78.4),('West Virginia',72.2),('Wisconsin',78.1),('Wyoming',76.8)
)
INSERT INTO metric_values(state_id,metric_id,year,value,source_record_id,import_id)
SELECT s.id,m.id,2022,o.value,o.state_name,i.id
FROM observations o
JOIN states s ON s.name=o.state_name
JOIN metrics m ON m.slug='life-expectancy'
JOIN imports i ON i.checksum='bundled-cdc-life-2022-v1';

WITH observations(state_name,value) AS (VALUES
('Alabama',88.823),('Alaska',102.359),('Arizona',100.677),('Arkansas',86.937),('California',110.720),
('Colorado',103.052),('Connecticut',103.610),('Delaware',99.808),('Florida',103.414),('Georgia',96.293),
('Hawaii',109.951),('Idaho',95.494),('Illinois',99.958),('Indiana',93.329),('Iowa',87.762),
('Kansas',90.068),('Kentucky',90.159),('Louisiana',88.207),('Maine',97.050),('Maryland',104.959),
('Massachusetts',105.757),('Michigan',96.217),('Minnesota',98.621),('Mississippi',86.953),('Missouri',90.817),
('Montana',94.645),('Nebraska',90.103),('Nevada',99.979),('New Hampshire',104.165),('New Jersey',108.805),
('New Mexico',92.212),('New York',107.921),('North Carolina',94.326),('North Dakota',88.959),('Ohio',92.774),
('Oklahoma',87.843),('Oregon',103.361),('Pennsylvania',97.572),('Rhode Island',102.280),('South Carolina',93.749),
('South Dakota',88.586),('Tennessee',91.870),('Texas',97.057),('Utah',98.864),('Vermont',97.958),
('Virginia',101.104),('Washington',107.013),('West Virginia',89.497),('Wisconsin',94.095),('Wyoming',92.691)
)
INSERT INTO metric_values(state_id,metric_id,year,value,source_record_id,import_id)
SELECT s.id,m.id,2024,o.value,o.state_name,i.id
FROM observations o
JOIN states s ON s.name=o.state_name
JOIN metrics m ON m.slug='cost-of-living-index'
JOIN imports i ON i.checksum='bundled-bea-rpp-2024-v1';

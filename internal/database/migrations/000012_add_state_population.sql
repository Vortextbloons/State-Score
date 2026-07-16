-- Add Census Vintage 2025 resident population estimates to state reference data.
ALTER TABLE states ADD COLUMN population INTEGER;
ALTER TABLE states ADD COLUMN population_year INTEGER;
ALTER TABLE states ADD COLUMN population_source_id INTEGER REFERENCES data_sources(id);

INSERT INTO data_sources(name,publisher,source_url,license,format,description) VALUES
('Census Vintage 2025 State Population Estimates','U.S. Census Bureau','https://www2.census.gov/programs-surveys/popest/datasets/2020-2025/state/totals/NST-EST2025-ALLDATA.csv','U.S. government public data','csv','Resident population estimates for July 1, 2025 from the Vintage 2025 Population Estimates Program.');

WITH p(state_name,population) AS (VALUES
('Alabama',5193088),('Alaska',737270),('Arizona',7623818),('Arkansas',3114791),('California',39355309),('Colorado',6012561),('Connecticut',3688496),('Delaware',1059952),('Florida',23462518),('Georgia',11302748),
('Hawaii',1432820),('Idaho',2029733),('Illinois',12719141),('Indiana',6973333),('Iowa',3238387),('Kansas',2977220),('Kentucky',4606864),('Louisiana',4618189),('Maine',1414874),('Maryland',6265347),
('Massachusetts',7154084),('Michigan',10127884),('Minnesota',5830405),('Mississippi',2954160),('Missouri',6270541),('Montana',1144694),('Nebraska',2018006),('Nevada',3282188),('New Hampshire',1415342),('New Jersey',9548215),
('New Mexico',2125498),('New York',20002427),('North Carolina',11197968),('North Dakota',799358),('Ohio',11900510),('Oklahoma',4123288),('Oregon',4273586),('Pennsylvania',13059432),('Rhode Island',1114521),('South Carolina',5570274),
('South Dakota',935094),('Tennessee',7315076),('Texas',31709821),('Utah',3538904),('Vermont',644663),('Virginia',8880107),('Washington',8001020),('West Virginia',1766147),('Wisconsin',5972787),('Wyoming',588753)
)
UPDATE states SET population=(SELECT p.population FROM p WHERE p.state_name=states.name), population_year=2025,
 population_source_id=(SELECT id FROM data_sources WHERE name='Census Vintage 2025 State Population Estimates'), updated_at=datetime('now')
WHERE name IN (SELECT state_name FROM p);

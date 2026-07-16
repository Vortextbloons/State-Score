-- Add five foundational metrics with complete, official 2024 state data.

INSERT INTO data_sources(name,publisher,source_url,license,format,description) VALUES
('ACS 2024 Labor Force Participation','U.S. Census Bureau','https://api.census.gov/data/2024/acs/acs1/subject','U.S. government public data','json','ACS 1-Year Subject Table S2301 labor-force participation rate for the population age 16 and over.'),
('NAEP 2024 Reading and Mathematics','National Center for Education Statistics','https://www.nationsreportcard.gov/reports/','U.S. government public data','html/json','Mean of 2024 public-school average scale scores for grade 4 and grade 8 reading and mathematics.'),
('ACS 2024 Health Insurance Coverage','U.S. Census Bureau','https://api.census.gov/data/2024/acs/acs1/subject','U.S. government public data','json','ACS 1-Year Subject Table S2701 percent of the civilian noninstitutionalized population without health insurance.'),
('NVSS 2024 Homicide Mortality','Centers for Disease Control and Prevention, National Center for Health Statistics','https://www.cdc.gov/nchs/state-stats/deaths/homicide.html','U.S. government public data','json/csv','CDC WONDER final 2024 age-adjusted homicide deaths per 100,000 population.'),
('ACS 2024 Owner Housing-Cost Burden','U.S. Census Bureau','https://api.census.gov/data/2024/acs/acs1','U.S. government public data','json','ACS 1-Year Detailed Table B25091 owner households spending at least 30 percent of income on selected monthly owner costs, excluding households for which the percentage cannot be computed.');

WITH v(category_slug,slug,name,description,unit,higher,source_name) AS (VALUES
 ('economy','labor-force-participation-rate','Labor-force participation rate','Percentage of the population age 16 and over in the civilian labor force.','Percent',1,'ACS 2024 Labor Force Participation'),
 ('education','naep-achievement-composite','NAEP achievement composite','Mean of grade 4 and grade 8 public-school average scale scores in reading and mathematics.','Scale score',1,'NAEP 2024 Reading and Mathematics'),
 ('safety','age-adjusted-homicide-death-rate','Age-adjusted homicide death rate','Homicide deaths per 100,000 population, age-adjusted to the 2000 U.S. standard population.','Per 100k',0,'NVSS 2024 Homicide Mortality'),
 ('affordability','owner-housing-cost-burden','Owner housing-cost burden','Percentage of owner households with computable housing-cost ratios spending at least 30 percent of income on selected monthly owner costs.','Percent',0,'ACS 2024 Owner Housing-Cost Burden')
)
INSERT INTO metrics(category_id,slug,name,description,unit,higher_is_better,normalization_method,default_weight,source_id,active)
SELECT c.id,v.slug,v.name,v.description,v.unit,v.higher,'percentile',1.0,ds.id,1
FROM v JOIN categories c ON c.slug=v.category_slug JOIN data_sources ds ON ds.name=v.source_name;

UPDATE metrics SET
 name='Uninsured rate',
 description='Percentage of the civilian noninstitutionalized population without health insurance coverage.',
 source_id=(SELECT id FROM data_sources WHERE name='ACS 2024 Health Insurance Coverage'),
 active=1,
 updated_at=datetime('now')
WHERE slug='uninsured-rate';

UPDATE metrics SET default_weight=1.0/(
 SELECT count(*) FROM metrics sibling WHERE sibling.category_id=metrics.category_id AND sibling.active=1
) WHERE active=1;

INSERT OR REPLACE INTO profile_metric_weights(profile_id,metric_id,weight)
SELECT p.id,m.id,m.default_weight FROM scoring_profiles p JOIN metrics m ON m.active=1;

INSERT INTO imports(source_id,status,started_at,completed_at,records_read,records_inserted,records_rejected,checksum)
SELECT id,'completed',datetime('now'),datetime('now'),50,50,0,'bundled-acs-foundational-2024-v1' FROM data_sources WHERE name='ACS 2024 Labor Force Participation';
INSERT INTO imports(source_id,status,started_at,completed_at,records_read,records_inserted,records_rejected,checksum)
SELECT id,'completed',datetime('now'),datetime('now'),50,50,0,'bundled-naep-composite-2024-v1' FROM data_sources WHERE name='NAEP 2024 Reading and Mathematics';
INSERT INTO imports(source_id,status,started_at,completed_at,records_read,records_inserted,records_rejected,checksum)
SELECT id,'completed',datetime('now'),datetime('now'),50,50,0,'bundled-acs-uninsured-2024-v1' FROM data_sources WHERE name='ACS 2024 Health Insurance Coverage';
INSERT INTO imports(source_id,status,started_at,completed_at,records_read,records_inserted,records_rejected,checksum)
SELECT id,'completed',datetime('now'),datetime('now'),50,50,0,'bundled-nvss-homicide-2024-v1' FROM data_sources WHERE name='NVSS 2024 Homicide Mortality';
INSERT INTO imports(source_id,status,started_at,completed_at,records_read,records_inserted,records_rejected,checksum)
SELECT id,'completed',datetime('now'),datetime('now'),50,50,0,'bundled-acs-owner-burden-2024-v1' FROM data_sources WHERE name='ACS 2024 Owner Housing-Cost Burden';

WITH o(state_name,participation,uninsured,owner_burden) AS (VALUES
('Alabama',59.2,8.2,11.4737),('Alaska',66.6,11.0,13.2771),('Arizona',61.4,10.3,14.7418),('Arkansas',59.0,9.4,11.1765),('California',64.2,5.9,20.1182),
('Colorado',68.7,7.9,17.4603),('Connecticut',66.1,5.8,17.3961),('Delaware',62.6,6.9,14.5865),('Florida',60.8,10.9,18.8506),('Georgia',64.5,12.0,14.3279),
('Hawaii',64.1,3.5,19.3544),('Idaho',63.2,9.2,13.5924),('Illinois',65.5,6.9,14.0193),('Indiana',64.2,7.5,11.8435),('Iowa',66.5,5.4,11.4286),
('Kansas',66.6,8.5,11.9232),('Kentucky',60.2,6.8,11.2191),('Louisiana',59.9,7.7,12.8872),('Maine',61.9,5.5,12.8952),('Maryland',67.9,6.3,15.0144),
('Massachusetts',67.2,2.8,17.0909),('Michigan',61.6,5.1,12.546),('Minnesota',68.1,5.1,13.2366),('Mississippi',58.9,9.7,11.8721),('Missouri',63.5,7.7,11.9589),
('Montana',63.2,8.8,15.6161),('Nebraska',68.7,7.1,12.7988),('Nevada',64.6,11.4,16.5862),('New Hampshire',66.4,4.5,16.1407),('New Jersey',66.9,7.7,18.316),
('New Mexico',58.7,10.1,12.0918),('New York',63.0,5.0,16.8945),('North Carolina',63.4,8.6,13.1919),('North Dakota',68.4,6.1,10.0825),('Ohio',63.5,6.7,11.9838),
('Oklahoma',61.5,11.5,12.2966),('Oregon',62.9,5.2,16.8771),('Pennsylvania',63.2,5.8,12.3305),('Rhode Island',65.6,4.6,19.0462),('South Carolina',61.3,9.0,12.736),
('South Dakota',67.9,8.1,12.5427),('Tennessee',63.3,9.7,12.0397),('Texas',65.7,16.7,15.1438),('Utah',69.8,8.3,14.7398),('Vermont',64.8,4.2,16.6202),
('Virginia',65.2,6.9,13.7143),('Washington',65.0,6.5,16.2249),('West Virginia',54.9,5.8,8.7666),('Wisconsin',65.4,5.3,12.1723),('Wyoming',64.9,10.3,14.2037)
), e(state_name,slug,value,checksum) AS (
 SELECT state_name,'labor-force-participation-rate',participation,'bundled-acs-foundational-2024-v1' FROM o UNION ALL
 SELECT state_name,'uninsured-rate',uninsured,'bundled-acs-uninsured-2024-v1' FROM o UNION ALL
 SELECT state_name,'owner-housing-cost-burden',owner_burden,'bundled-acs-owner-burden-2024-v1' FROM o
)
INSERT INTO metric_values(state_id,metric_id,year,value,source_record_id,import_id)
SELECT s.id,m.id,2024,e.value,e.state_name,i.id FROM e JOIN states s ON s.name=e.state_name JOIN metrics m ON m.slug=e.slug JOIN imports i ON i.checksum=e.checksum;

WITH o(state_name,value) AS (VALUES
('Alabama',240.25),('Alaska',234.50),('Arizona',241.00),('Arkansas',240.25),('California',242.00),('Colorado',250.75),('Connecticut',249.50),('Delaware',238.75),('Florida',245.25),('Georgia',244.50),
('Hawaii',245.50),('Idaho',248.25),('Illinois',247.25),('Indiana',250.00),('Iowa',247.00),('Kansas',244.75),('Kentucky',246.25),('Louisiana',243.75),('Maine',242.75),('Maryland',244.00),
('Massachusetts',255.50),('Michigan',242.25),('Minnesota',249.25),('Mississippi',245.00),('Missouri',243.00),('Montana',248.75),('Nebraska',246.50),('Nevada',241.00),('New Hampshire',251.75),('New Jersey',252.50),
('New Mexico',231.50),('New York',244.25),('North Carolina',245.75),('North Dakota',248.50),('Ohio',248.50),('Oklahoma',238.25),('Oregon',239.75),('Pennsylvania',247.25),('Rhode Island',245.25),('South Carolina',243.75),
('South Dakota',248.75),('Tennessee',247.50),('Texas',243.50),('Utah',251.00),('Vermont',245.25),('Virginia',245.75),('Washington',246.75),('West Virginia',236.50),('Wisconsin',249.25),('Wyoming',250.75)
)
INSERT INTO metric_values(state_id,metric_id,year,value,source_record_id,import_id)
SELECT s.id,m.id,2024,o.value,o.state_name,i.id FROM o JOIN states s ON s.name=o.state_name JOIN metrics m ON m.slug='naep-achievement-composite' JOIN imports i ON i.checksum='bundled-naep-composite-2024-v1';

WITH o(state_name,value) AS (VALUES
('Alabama',13.3),('Alaska',7.7),('Arizona',6.4),('Arkansas',8.5),('California',4.5),('Colorado',5.1),('Connecticut',3.0),('Delaware',6.9),('Florida',5.5),('Georgia',9.2),
('Hawaii',2.8),('Idaho',1.8),('Illinois',8.7),('Indiana',7.3),('Iowa',3.4),('Kansas',5.6),('Kentucky',7.0),('Louisiana',15.2),('Maine',2.6),('Maryland',9.0),
('Massachusetts',2.2),('Michigan',5.7),('Minnesota',3.4),('Mississippi',19.7),('Missouri',9.4),('Montana',4.5),('Nebraska',3.1),('Nevada',7.0),('New Hampshire',1.4),('New Jersey',2.9),
('New Mexico',13.2),('New York',3.5),('North Carolina',7.8),('North Dakota',2.5),('Ohio',6.6),('Oklahoma',7.2),('Oregon',4.0),('Pennsylvania',5.6),('Rhode Island',2.0),('South Carolina',9.9),
('South Dakota',6.6),('Tennessee',10.0),('Texas',5.9),('Utah',2.3),('Vermont',3.2),('Virginia',5.6),('Washington',4.4),('West Virginia',5.9),('Wisconsin',4.8),('Wyoming',4.8)
)
INSERT INTO metric_values(state_id,metric_id,year,value,source_record_id,import_id)
SELECT s.id,m.id,2024,o.value,o.state_name,i.id FROM o JOIN states s ON s.name=o.state_name JOIN metrics m ON m.slug='age-adjusted-homicide-death-rate' JOIN imports i ON i.checksum='bundled-nvss-homicide-2024-v1';

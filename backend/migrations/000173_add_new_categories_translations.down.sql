-- Remove translations for categories added in migration 000172
DELETE FROM translations 
WHERE entity_type = 'category' 
AND entity_id IN (2100, 2101, 2102, 2103, 2200, 2201, 2202, 2203, 2300, 2301, 2302, 2303, 2400, 2401, 2402, 2500, 2501, 2502);
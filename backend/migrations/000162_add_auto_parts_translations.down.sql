-- Remove translations for auto parts categories

DELETE FROM translations 
WHERE entity_type = 'category' 
AND entity_id IN (
    1304, 1305, 1306, 1307, 1308, 1309, 1310, 1311, 1312, 1313,
    1314, 1315, 1316, 1317, 1318, 1319, 1320, 1321, 1322, 1323,
    1324, 1325, 1326, 1327, 1328, 1329, 1330, 1331, 1332, 1333,
    1334, 1335, 1336, 1337, 1338, 1339, 1340
);

-- Restore original translations for category 1303 if they were different
UPDATE translations 
SET translated_text = 'Auto delovi' 
WHERE entity_type = 'category' AND entity_id = 1303 AND field_name = 'name' AND language = 'en';

UPDATE translations 
SET translated_text = 'Auto delovi' 
WHERE entity_type = 'category' AND entity_id = 1303 AND field_name = 'name' AND language = 'ru';
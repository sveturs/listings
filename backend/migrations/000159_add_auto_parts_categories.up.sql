-- Add automotive parts categories with proper hierarchy

-- First level - under Auto delovi (id=1303)
INSERT INTO marketplace_categories (id, parent_id, name, slug, level, sort_order, is_active) VALUES
-- Gume i točkovi (Шины и колеса)
(1304, 1303, 'Gume i točkovi', 'tires-and-wheels', 1, 1, true),
-- Motor i delovi motora
(1305, 1303, 'Motor i delovi motora', 'engine-and-parts', 1, 2, true),
-- Karoserija i delovi
(1306, 1303, 'Karoserija i delovi', 'body-parts', 1, 3, true),
-- Električni i elektronski delovi
(1307, 1303, 'Električni i elektronski delovi', 'electrical-parts', 1, 4, true),
-- Sistem za kočenje
(1308, 1303, 'Sistem za kočenje', 'brake-system', 1, 5, true),
-- Sistem vešanja
(1309, 1303, 'Sistem vešanja', 'suspension-system', 1, 6, true),
-- Sistem hlađenja
(1310, 1303, 'Sistem hlađenja', 'cooling-system', 1, 7, true),
-- Transmisija i delovi
(1311, 1303, 'Transmisija i delovi', 'transmission-parts', 1, 8, true),
-- Unutrašnjost
(1312, 1303, 'Unutrašnjost', 'interior-parts', 1, 9, true),
-- Dodatna oprema
(1313, 1303, 'Dodatna oprema', 'auto-accessories', 1, 10, true);

-- Second level - under Gume i točkovi (id=1304)
INSERT INTO marketplace_categories (id, parent_id, name, slug, level, sort_order, is_active) VALUES
-- Letnje gume
(1314, 1304, 'Letnje gume', 'summer-tires', 2, 1, true),
-- Zimske gume', 
(1315, 1304, 'Zimske gume', 'winter-tires', 2, 2, true),
-- Celogodišnje gume
(1316, 1304, 'Celogodišnje gume', 'all-season-tires', 2, 3, true),
-- Felne
(1317, 1304, 'Felne', 'rims', 2, 4, true),
-- Kompletni točkovi
(1318, 1304, 'Kompletni točkovi', 'complete-wheels', 2, 5, true),
-- Ratkapne
(1319, 1304, 'Ratkapne', 'wheel-covers', 2, 6, true),
-- Vijci za točkove
(1320, 1304, 'Vijci za točkove', 'wheel-bolts', 2, 7, true);

-- Third level - under different tire types
INSERT INTO marketplace_categories (id, parent_id, name, slug, level, sort_order, is_active) VALUES
-- Under Letnje gume (1314)
(1321, 1314, 'Putničke letnje gume', 'passenger-summer-tires', 3, 1, true),
(1322, 1314, 'SUV letnje gume', 'suv-summer-tires', 3, 2, true),
(1323, 1314, 'Kamionske letnje gume', 'truck-summer-tires', 3, 3, true),

-- Under Zimske gume (1315)
(1324, 1315, 'Putničke zimske gume', 'passenger-winter-tires', 3, 1, true),
(1325, 1315, 'SUV zimske gume', 'suv-winter-tires', 3, 2, true),
(1326, 1315, 'Kamionske zimske gume', 'truck-winter-tires', 3, 3, true),

-- Under Felne (1317)
(1327, 1317, 'Čelične felne', 'steel-rims', 3, 1, true),
(1328, 1317, 'Aluminijumske felne', 'aluminum-rims', 3, 2, true),
(1329, 1317, 'Sportske felne', 'sport-rims', 3, 3, true);

-- Second level - under Motor i delovi motora (id=1305)
INSERT INTO marketplace_categories (id, parent_id, name, slug, level, sort_order, is_active) VALUES
(1330, 1305, 'Filtri', 'filters', 2, 1, true),
(1331, 1305, 'Remeni i lančanici', 'belts-and-chains', 2, 2, true),
(1332, 1305, 'Ulje i tečnosti', 'oils-and-fluids', 2, 3, true),
(1333, 1305, 'Svećice', 'spark-plugs', 2, 4, true),
(1334, 1305, 'Izduvni sistem', 'exhaust-system', 2, 5, true);

-- Second level - under Karoserija i delovi (id=1306)
INSERT INTO marketplace_categories (id, parent_id, name, slug, level, sort_order, is_active) VALUES
(1335, 1306, 'Branici', 'bumpers', 2, 1, true),
(1336, 1306, 'Vrata', 'doors', 2, 2, true),
(1337, 1306, 'Haube', 'hoods', 2, 3, true),
(1338, 1306, 'Blatobrani', 'fenders', 2, 4, true),
(1339, 1306, 'Retrovizori', 'mirrors', 2, 5, true),
(1340, 1306, 'Stakla', 'windows', 2, 6, true);

-- Update parent category level and make sure hierarchy is correct
UPDATE marketplace_categories 
SET level = 1 
WHERE parent_id = 1303;

UPDATE marketplace_categories 
SET level = 2 
WHERE parent_id IN (SELECT id FROM marketplace_categories WHERE parent_id = 1303);

UPDATE marketplace_categories 
SET level = 3 
WHERE parent_id IN (SELECT id FROM marketplace_categories WHERE level = 2);

-- Fix parent category level 
UPDATE marketplace_categories
SET level = 1
WHERE id = 1303;

UPDATE marketplace_categories
SET level = 0
WHERE id = 1003;
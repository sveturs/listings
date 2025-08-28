-- Update Serbian translations from Cyrillic to Latin for auto parts categories

UPDATE translations 
SET translated_text = CASE 
    -- Main categories
    WHEN entity_id = 1304 THEN 'Gume i točkovi'
    WHEN entity_id = 1305 THEN 'Motor i delovi motora'
    WHEN entity_id = 1306 THEN 'Karoserija i delovi'
    WHEN entity_id = 1307 THEN 'Električni i elektronski delovi'
    WHEN entity_id = 1308 THEN 'Sistem za kočenje'
    WHEN entity_id = 1309 THEN 'Sistem vešanja'
    WHEN entity_id = 1310 THEN 'Sistem hlađenja'
    WHEN entity_id = 1311 THEN 'Transmisija i delovi'
    WHEN entity_id = 1312 THEN 'Unutrašnjost'
    WHEN entity_id = 1313 THEN 'Dodatna oprema'
    
    -- Tire subcategories
    WHEN entity_id = 1314 THEN 'Letnje gume'
    WHEN entity_id = 1315 THEN 'Zimske gume'
    WHEN entity_id = 1316 THEN 'Celogodišnje gume'
    WHEN entity_id = 1317 THEN 'Felne'
    WHEN entity_id = 1318 THEN 'Kompletni točkovi'
    WHEN entity_id = 1319 THEN 'Ratkapne'
    WHEN entity_id = 1320 THEN 'Vijci za točkove'
    
    -- Third level tire categories
    WHEN entity_id = 1321 THEN 'Putničke letnje gume'
    WHEN entity_id = 1322 THEN 'SUV letnje gume'
    WHEN entity_id = 1323 THEN 'Kamionske letnje gume'
    WHEN entity_id = 1324 THEN 'Putničke zimske gume'
    WHEN entity_id = 1325 THEN 'SUV zimske gume'
    WHEN entity_id = 1326 THEN 'Kamionske zimske gume'
    WHEN entity_id = 1327 THEN 'Čelične felne'
    WHEN entity_id = 1328 THEN 'Aluminijumske felne'
    WHEN entity_id = 1329 THEN 'Sportske felne'
    
    -- Engine subcategories
    WHEN entity_id = 1330 THEN 'Filtri'
    WHEN entity_id = 1331 THEN 'Remeni i lančanici'
    WHEN entity_id = 1332 THEN 'Ulje i tečnosti'
    WHEN entity_id = 1333 THEN 'Svećice'
    WHEN entity_id = 1334 THEN 'Izduvni sistem'
    
    -- Body parts subcategories
    WHEN entity_id = 1335 THEN 'Branici'
    WHEN entity_id = 1336 THEN 'Vrata'
    WHEN entity_id = 1337 THEN 'Haube'
    WHEN entity_id = 1338 THEN 'Blatobrani'
    WHEN entity_id = 1339 THEN 'Retrovizori'
    WHEN entity_id = 1340 THEN 'Stakla'
END
WHERE entity_type = 'category' 
AND entity_id IN (
    1304, 1305, 1306, 1307, 1308, 1309, 1310, 1311, 1312, 1313,
    1314, 1315, 1316, 1317, 1318, 1319, 1320, 1321, 1322, 1323,
    1324, 1325, 1326, 1327, 1328, 1329, 1330, 1331, 1332, 1333,
    1334, 1335, 1336, 1337, 1338, 1339, 1340
)
AND field_name = 'name'
AND language = 'sr';

-- Update SEO titles and descriptions too
UPDATE translations 
SET translated_text = CASE 
    WHEN entity_id = 1304 AND field_name = 'seo_title' THEN 'Gume, točkovi i felne'
    WHEN entity_id = 1304 AND field_name = 'seo_description' THEN 'Letnje, zimske i celogodišnje gume, točkovi i felne'
    WHEN entity_id = 1305 AND field_name = 'seo_title' THEN 'Delovi i komponente motora'
    WHEN entity_id = 1306 AND field_name = 'seo_title' THEN 'Delovi karoserije i pribor'
END
WHERE entity_type = 'category' 
AND entity_id IN (1304, 1305, 1306)
AND field_name IN ('seo_title', 'seo_description')
AND language = 'sr';

-- Update tire attribute translations
UPDATE translations 
SET translated_text = CASE 
    WHEN field_name = 'display_name' AND entity_id = (SELECT id FROM category_attributes WHERE name = 'tire_width') THEN 'Širina gume'
    WHEN field_name = 'display_name' AND entity_id = (SELECT id FROM category_attributes WHERE name = 'tire_profile') THEN 'Profil gume'
    WHEN field_name = 'display_name' AND entity_id = (SELECT id FROM category_attributes WHERE name = 'tire_diameter') THEN 'Prečnik'
    WHEN field_name = 'display_name' AND entity_id = (SELECT id FROM category_attributes WHERE name = 'tire_season') THEN 'Sezona'
    WHEN field_name = 'display_name' AND entity_id = (SELECT id FROM category_attributes WHERE name = 'tire_brand') THEN 'Proizvođač'
    WHEN field_name = 'display_name' AND entity_id = (SELECT id FROM category_attributes WHERE name = 'tire_condition') THEN 'Stanje'
    WHEN field_name = 'display_name' AND entity_id = (SELECT id FROM category_attributes WHERE name = 'tread_depth') THEN 'Dubina šare'
    WHEN field_name = 'display_name' AND entity_id = (SELECT id FROM category_attributes WHERE name = 'tire_year') THEN 'Godina proizvodnje'
    WHEN field_name = 'display_name' AND entity_id = (SELECT id FROM category_attributes WHERE name = 'tire_quantity') THEN 'Količina'
END
WHERE entity_type = 'attribute' 
AND language = 'sr'
AND entity_id IN (
    SELECT id FROM category_attributes 
    WHERE name IN ('tire_width', 'tire_profile', 'tire_diameter', 'tire_season', 
                   'tire_brand', 'tire_condition', 'tread_depth', 'tire_year', 'tire_quantity')
);

-- Update attribute option translations
UPDATE attribute_option_translations 
SET sr_translation = CASE 
    -- Season translations
    WHEN attribute_name = 'tire_season' AND option_value = 'summer' THEN 'Letnje'
    WHEN attribute_name = 'tire_season' AND option_value = 'winter' THEN 'Zimske'
    WHEN attribute_name = 'tire_season' AND option_value = 'all-season' THEN 'Celogodišnje'
    
    -- Condition translations
    WHEN attribute_name = 'tire_condition' AND option_value = 'new' THEN 'Nove'
    WHEN attribute_name = 'tire_condition' AND option_value = 'used' THEN 'Polovne'
    
    -- Quantity translations
    WHEN attribute_name = 'tire_quantity' AND option_value = '1' THEN '1 kom'
    WHEN attribute_name = 'tire_quantity' AND option_value = '2' THEN '2 kom'
    WHEN attribute_name = 'tire_quantity' AND option_value = '3' THEN '3 kom'
    WHEN attribute_name = 'tire_quantity' AND option_value = '4' THEN '4 kom'
    WHEN attribute_name = 'tire_quantity' AND option_value = 'set' THEN 'Komplet (4 kom)'
END
WHERE attribute_name IN ('tire_season', 'tire_condition', 'tire_quantity');
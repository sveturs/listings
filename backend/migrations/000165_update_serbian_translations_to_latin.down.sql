-- Revert Serbian translations back to Cyrillic for auto parts categories

UPDATE translations 
SET translated_text = CASE 
    -- Main categories
    WHEN entity_id = 1304 THEN 'Гуме и точкови'
    WHEN entity_id = 1305 THEN 'Мотор и делови мотора'
    WHEN entity_id = 1306 THEN 'Каросерија и делови'
    WHEN entity_id = 1307 THEN 'Електрични и електронски делови'
    WHEN entity_id = 1308 THEN 'Систем за кочење'
    WHEN entity_id = 1309 THEN 'Систем вешања'
    WHEN entity_id = 1310 THEN 'Систем хлађења'
    WHEN entity_id = 1311 THEN 'Трансмисија и делови'
    WHEN entity_id = 1312 THEN 'Унутрашњост'
    WHEN entity_id = 1313 THEN 'Додатна опрема'
    
    -- Tire subcategories
    WHEN entity_id = 1314 THEN 'Летње гуме'
    WHEN entity_id = 1315 THEN 'Зимске гуме'
    WHEN entity_id = 1316 THEN 'Целогодишње гуме'
    WHEN entity_id = 1317 THEN 'Фелне'
    WHEN entity_id = 1318 THEN 'Комплетни точкови'
    WHEN entity_id = 1319 THEN 'Раткапне'
    WHEN entity_id = 1320 THEN 'Вијци за точкове'
    
    -- Third level tire categories
    WHEN entity_id = 1321 THEN 'Путничке летње гуме'
    WHEN entity_id = 1322 THEN 'SUV летње гуме'
    WHEN entity_id = 1323 THEN 'Камионске летње гуме'
    WHEN entity_id = 1324 THEN 'Путничке зимске гуме'
    WHEN entity_id = 1325 THEN 'SUV зимске гуме'
    WHEN entity_id = 1326 THEN 'Камионске зимске гуме'
    WHEN entity_id = 1327 THEN 'Челичне фелне'
    WHEN entity_id = 1328 THEN 'Алуминијумске фелне'
    WHEN entity_id = 1329 THEN 'Спортске фелне'
    
    -- Engine subcategories
    WHEN entity_id = 1330 THEN 'Филтри'
    WHEN entity_id = 1331 THEN 'Ремени и ланчаници'
    WHEN entity_id = 1332 THEN 'Уље и течности'
    WHEN entity_id = 1333 THEN 'Свећице'
    WHEN entity_id = 1334 THEN 'Издувни систем'
    
    -- Body parts subcategories
    WHEN entity_id = 1335 THEN 'Браници'
    WHEN entity_id = 1336 THEN 'Врата'
    WHEN entity_id = 1337 THEN 'Хаубе'
    WHEN entity_id = 1338 THEN 'Блатобрани'
    WHEN entity_id = 1339 THEN 'Ретровизори'
    WHEN entity_id = 1340 THEN 'Стакла'
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

-- Revert SEO titles and descriptions
UPDATE translations 
SET translated_text = CASE 
    WHEN entity_id = 1304 AND field_name = 'seo_title' THEN 'Гуме, точкови и фелне'
    WHEN entity_id = 1304 AND field_name = 'seo_description' THEN 'Летње, зимске и целогодишње гуме, точкови и фелне'
    WHEN entity_id = 1305 AND field_name = 'seo_title' THEN 'Делови и компоненте мотора'
    WHEN entity_id = 1306 AND field_name = 'seo_title' THEN 'Делови каросерије и прибор'
END
WHERE entity_type = 'category' 
AND entity_id IN (1304, 1305, 1306)
AND field_name IN ('seo_title', 'seo_description')
AND language = 'sr';

-- Revert tire attribute translations
UPDATE translations 
SET translated_text = CASE 
    WHEN field_name = 'display_name' AND entity_id = (SELECT id FROM category_attributes WHERE name = 'tire_width') THEN 'Ширина гуме'
    WHEN field_name = 'display_name' AND entity_id = (SELECT id FROM category_attributes WHERE name = 'tire_profile') THEN 'Профил гуме'
    WHEN field_name = 'display_name' AND entity_id = (SELECT id FROM category_attributes WHERE name = 'tire_diameter') THEN 'Пречник'
    WHEN field_name = 'display_name' AND entity_id = (SELECT id FROM category_attributes WHERE name = 'tire_season') THEN 'Сезона'
    WHEN field_name = 'display_name' AND entity_id = (SELECT id FROM category_attributes WHERE name = 'tire_brand') THEN 'Произвођач'
    WHEN field_name = 'display_name' AND entity_id = (SELECT id FROM category_attributes WHERE name = 'tire_condition') THEN 'Стање'
    WHEN field_name = 'display_name' AND entity_id = (SELECT id FROM category_attributes WHERE name = 'tread_depth') THEN 'Дубина шаре'
    WHEN field_name = 'display_name' AND entity_id = (SELECT id FROM category_attributes WHERE name = 'tire_year') THEN 'Година производње'
    WHEN field_name = 'display_name' AND entity_id = (SELECT id FROM category_attributes WHERE name = 'tire_quantity') THEN 'Количина'
END
WHERE entity_type = 'attribute' 
AND language = 'sr'
AND entity_id IN (
    SELECT id FROM category_attributes 
    WHERE name IN ('tire_width', 'tire_profile', 'tire_diameter', 'tire_season', 
                   'tire_brand', 'tire_condition', 'tread_depth', 'tire_year', 'tire_quantity')
);

-- Revert attribute option translations
UPDATE attribute_option_translations 
SET sr_translation = CASE 
    -- Season translations
    WHEN attribute_name = 'tire_season' AND option_value = 'summer' THEN 'Летње'
    WHEN attribute_name = 'tire_season' AND option_value = 'winter' THEN 'Зимске'
    WHEN attribute_name = 'tire_season' AND option_value = 'all-season' THEN 'Целогодишње'
    
    -- Condition translations
    WHEN attribute_name = 'tire_condition' AND option_value = 'new' THEN 'Нове'
    WHEN attribute_name = 'tire_condition' AND option_value = 'used' THEN 'Половне'
    
    -- Quantity translations
    WHEN attribute_name = 'tire_quantity' AND option_value = '1' THEN '1 ком'
    WHEN attribute_name = 'tire_quantity' AND option_value = '2' THEN '2 ком'
    WHEN attribute_name = 'tire_quantity' AND option_value = '3' THEN '3 ком'
    WHEN attribute_name = 'tire_quantity' AND option_value = '4' THEN '4 ком'
    WHEN attribute_name = 'tire_quantity' AND option_value = 'set' THEN 'Комплет (4 ком)'
END
WHERE attribute_name IN ('tire_season', 'tire_condition', 'tire_quantity');
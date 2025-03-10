-- backend/migrations/000014_add_attribute_translations.up.sql
-- Добавляем переводы для атрибутов на английский
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at) 
SELECT 'attribute', id, 'en', 'display_name', 
    CASE 
        -- Автомобили
        WHEN name = 'make' THEN 'Make'
        WHEN name = 'model' THEN 'Model'
        WHEN name = 'year' THEN 'Year'
        WHEN name = 'mileage' THEN 'Mileage (km)'
        WHEN name = 'engine_capacity' THEN 'Engine capacity (L)'
        WHEN name = 'fuel_type' THEN 'Fuel type'
        WHEN name = 'transmission' THEN 'Transmission'
        WHEN name = 'body_type' THEN 'Body type'
        WHEN name = 'color' THEN 'Color'
        -- Недвижимость
        WHEN name = 'property_type' THEN 'Property type'
        WHEN name = 'rooms' THEN 'Rooms'
        WHEN name = 'floor' THEN 'Floor'
        WHEN name = 'total_floors' THEN 'Total floors'
        WHEN name = 'area' THEN 'Area (m²)'
        WHEN name = 'land_area' THEN 'Land area'
        WHEN name = 'building_type' THEN 'Building type'
        WHEN name = 'has_balcony' THEN 'Balcony'
        WHEN name = 'has_elevator' THEN 'Elevator'
        WHEN name = 'has_parking' THEN 'Parking'
        -- Телефоны
        WHEN name = 'brand' THEN 'Brand'
        WHEN name = 'model_phone' THEN 'Model'
        WHEN name = 'memory' THEN 'Memory (GB)'
        WHEN name = 'ram' THEN 'RAM (GB)'
        WHEN name = 'os' THEN 'Operating system'
        WHEN name = 'screen_size' THEN 'Screen size (inches)'
        WHEN name = 'camera' THEN 'Camera (MP)'
        WHEN name = 'has_5g' THEN '5G'
        -- Компьютеры
        WHEN name = 'pc_brand' THEN 'Brand'
        WHEN name = 'pc_type' THEN 'Type'
        WHEN name = 'cpu' THEN 'Processor'
        WHEN name = 'gpu' THEN 'Graphics card'
        WHEN name = 'ram_pc' THEN 'RAM (GB)'
        WHEN name = 'storage_type' THEN 'Storage type'
        WHEN name = 'storage_capacity' THEN 'Storage capacity (GB)'
        WHEN name = 'os_pc' THEN 'Operating system'
        ELSE display_name
    END, 
    false, true, NOW(), NOW()
FROM category_attributes
ON CONFLICT DO NOTHING;

-- Добавляем переводы для атрибутов на сербский
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at) 
SELECT 'attribute', id, 'sr', 'display_name', 
    CASE 
        -- Автомобили
        WHEN name = 'make' THEN 'Marka'
        WHEN name = 'model' THEN 'Model'
        WHEN name = 'year' THEN 'Godina proizvodnje'
        WHEN name = 'mileage' THEN 'Kilometraža'
        WHEN name = 'engine_capacity' THEN 'Zapremina motora (L)'
        WHEN name = 'fuel_type' THEN 'Vrsta goriva'
        WHEN name = 'transmission' THEN 'Menjač'
        WHEN name = 'body_type' THEN 'Tip karoserije'
        WHEN name = 'color' THEN 'Boja'
        -- Недвижимость
        WHEN name = 'property_type' THEN 'Tip nekretnine'
        WHEN name = 'rooms' THEN 'Broj soba'
        WHEN name = 'floor' THEN 'Sprat'
        WHEN name = 'total_floors' THEN 'Ukupno spratova'
        WHEN name = 'area' THEN 'Površina (m²)'
        WHEN name = 'land_area' THEN 'Površina zemljišta'
        WHEN name = 'building_type' THEN 'Tip zgrade'
        WHEN name = 'has_balcony' THEN 'Balkon'
        WHEN name = 'has_elevator' THEN 'Lift'
        WHEN name = 'has_parking' THEN 'Parking'
        -- Телефоны
        WHEN name = 'brand' THEN 'Brend'
        WHEN name = 'model_phone' THEN 'Model'
        WHEN name = 'memory' THEN 'Memorija (GB)'
        WHEN name = 'ram' THEN 'RAM (GB)'
        WHEN name = 'os' THEN 'Operativni sistem'
        WHEN name = 'screen_size' THEN 'Veličina ekrana (inči)'
        WHEN name = 'camera' THEN 'Kamera (MP)'
        WHEN name = 'has_5g' THEN '5G'
        -- Компьютеры
        WHEN name = 'pc_brand' THEN 'Brend'
        WHEN name = 'pc_type' THEN 'Tip'
        WHEN name = 'cpu' THEN 'Procesor'
        WHEN name = 'gpu' THEN 'Grafička kartica'
        WHEN name = 'ram_pc' THEN 'RAM (GB)'
        WHEN name = 'storage_type' THEN 'Tip skladišta'
        WHEN name = 'storage_capacity' THEN 'Kapacitet skladišta (GB)'
        WHEN name = 'os_pc' THEN 'Operativni sistem'
        ELSE display_name
    END, 
    false, true, NOW(), NOW()
FROM category_attributes
ON CONFLICT DO NOTHING;
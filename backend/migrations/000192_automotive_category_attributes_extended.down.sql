-- Индекс idx_listing_attributes_vin не создавался в up миграции

-- Удаляем переводы для опций оборудования
DELETE FROM attribute_option_translations 
WHERE attribute_name IN ('car_condition', 'equipment_features')
AND option_value IN (
    'new', 'used', 'damaged', 'for_parts',
    'abs', 'esp', 'airbag', 'climate_control', 'cruise_control',
    'parking_sensors', 'rear_camera', 'navigation', 'leather_seats',
    'heated_seats', 'sunroof', 'xenon', 'led', 'alloy_wheels', 'tow_bar'
);

-- Удаляем переводы атрибутов
DELETE FROM translations 
WHERE entity_type = 'attribute' 
AND entity_id IN (
    SELECT id FROM category_attributes 
    WHERE name IN (
        'vin_number', 'car_condition', 'owner_count', 'service_book', 
        'warranty', 'warranty_period', 'exchange_possible', 
        'financing_available', 'country_origin', 'first_registration',
        'inspection_valid_until', 'registration_valid_until',
        'equipment_features', 'additional_equipment'
    )
);

-- Удаляем маппинг атрибутов с категорией
DELETE FROM category_attribute_mapping 
WHERE category_id = 1301 
AND attribute_id IN (
    SELECT id FROM category_attributes 
    WHERE name IN (
        'vin_number', 'car_condition', 'owner_count', 'service_book', 
        'warranty', 'warranty_period', 'exchange_possible', 
        'financing_available', 'country_origin', 'first_registration',
        'inspection_valid_until', 'registration_valid_until',
        'equipment_features', 'additional_equipment'
    )
);

-- Удаляем новые атрибуты
DELETE FROM category_attributes 
WHERE name IN (
    'vin_number', 'car_condition', 'owner_count', 'service_book', 
    'warranty', 'warranty_period', 'exchange_possible', 
    'financing_available', 'country_origin', 'first_registration',
    'inspection_valid_until', 'registration_valid_until',
    'equipment_features', 'additional_equipment'
);
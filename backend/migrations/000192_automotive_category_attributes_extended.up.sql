-- Расширенные атрибуты для автомобильных категорий

-- Добавляем новые атрибуты (с проверкой на существование)
INSERT INTO category_attributes (name, display_name, attribute_type, options, is_required, sort_order) 
SELECT * FROM (VALUES
-- VIN номер
('vin_number', 'VIN Number', 'text', '{"placeholder": "VIN номер", "validation": {"pattern": "^[A-HJ-NPR-Z0-9]{17}$"}}'::jsonb, false, 150),

-- Состояние автомобиля
('car_condition', 'Car Condition', 'select', '{"options": [{"value": "new", "label": "New"}, {"value": "used", "label": "Used"}, {"value": "damaged", "label": "Damaged"}, {"value": "for_parts", "label": "For Parts"}]}'::jsonb, true, 160),

-- История владельцев
('owner_count', 'Number of Owners', 'number', '{"min": 1, "max": 10, "default": 1}'::jsonb, false, 170),

-- Сервисная книжка
('service_book', 'Service Book', 'boolean', '{"default": false}'::jsonb, false, 180),

-- Гарантия
('warranty', 'Warranty', 'boolean', '{"default": false}'::jsonb, false, 190),

-- Период гарантии
('warranty_period', 'Warranty Period', 'text', '{"placeholder": "Warranty period"}'::jsonb, false, 200),

-- Возможность обмена
('exchange_possible', 'Exchange Possible', 'boolean', '{"default": false}'::jsonb, false, 210),

-- Лизинг/Кредит
('financing_available', 'Financing Available', 'boolean', '{"default": false}'::jsonb, false, 220),

-- Страна происхождения
('country_origin', 'Country of Origin', 'text', '{"placeholder": "Country of origin"}'::jsonb, false, 230),

-- Дата первой регистрации
('first_registration', 'First Registration Date', 'date', '{}'::jsonb, false, 240),

-- Техосмотр до
('inspection_valid_until', 'Inspection Valid Until', 'date', '{}'::jsonb, false, 250),

-- Регистрация до
('registration_valid_until', 'Registration Valid Until', 'date', '{}'::jsonb, false, 260),

-- Особенности комплектации
('equipment_features', 'Equipment Features', 'multiselect', '{"options": [{"value": "abs", "label": "ABS"}, {"value": "esp", "label": "ESP"}, {"value": "asr", "label": "ASR"}, {"value": "airbag", "label": "Airbag"}, {"value": "climate_control", "label": "Climate Control"}, {"value": "cruise_control", "label": "Cruise Control"}, {"value": "parking_sensors", "label": "Parking Sensors"}, {"value": "rear_camera", "label": "Rear Camera"}, {"value": "navigation", "label": "Navigation"}, {"value": "leather_seats", "label": "Leather Seats"}, {"value": "heated_seats", "label": "Heated Seats"}, {"value": "sunroof", "label": "Sunroof"}, {"value": "xenon", "label": "Xenon"}, {"value": "led", "label": "LED"}, {"value": "alloy_wheels", "label": "Alloy Wheels"}, {"value": "tow_bar", "label": "Tow Bar"}]}'::jsonb, false, 270),

-- Дополнительное оборудование
('additional_equipment', 'Additional Equipment', 'textarea', '{"placeholder": "Additional equipment", "rows": 4}'::jsonb, false, 280)
) AS v(name, display_name, attribute_type, options, is_required, sort_order)
WHERE NOT EXISTS (SELECT 1 FROM category_attributes WHERE category_attributes.name = v.name);

-- Привязываем атрибуты к категории "Личные автомобили" (1301)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order) 
SELECT 1301, ca.id, true, ca.is_required, ca.sort_order
FROM category_attributes ca
WHERE ca.name IN (
    'vin_number', 'car_condition', 'owner_count', 'service_book', 
    'warranty', 'warranty_period', 'exchange_possible', 
    'financing_available', 'country_origin', 'first_registration',
    'inspection_valid_until', 'registration_valid_until',
    'equipment_features', 'additional_equipment'
)
AND NOT EXISTS (
    SELECT 1 FROM category_attribute_mapping cam 
    WHERE cam.category_id = 1301 AND cam.attribute_id = ca.id
);

-- Добавляем переводы для новых атрибутов (с проверкой на существование)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_verified)
SELECT * FROM (VALUES
-- VIN номер
('attribute', (SELECT id FROM category_attributes WHERE name = 'vin_number'), 'ru', 'label', 'VIN номер', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'vin_number'), 'sr', 'label', 'VIN broj', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'vin_number'), 'en', 'label', 'VIN number', true),

-- Состояние автомобиля
('attribute', (SELECT id FROM category_attributes WHERE name = 'car_condition'), 'ru', 'label', 'Состояние', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'car_condition'), 'sr', 'label', 'Stanje', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'car_condition'), 'en', 'label', 'Condition', true),

-- История владельцев
('attribute', (SELECT id FROM category_attributes WHERE name = 'owner_count'), 'ru', 'label', 'Количество владельцев', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'owner_count'), 'sr', 'label', 'Broj vlasnika', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'owner_count'), 'en', 'label', 'Number of owners', true),

-- Сервисная книжка
('attribute', (SELECT id FROM category_attributes WHERE name = 'service_book'), 'ru', 'label', 'Сервисная книжка', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'service_book'), 'sr', 'label', 'Servisna knjižica', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'service_book'), 'en', 'label', 'Service book', true),

-- Гарантия
('attribute', (SELECT id FROM category_attributes WHERE name = 'warranty'), 'ru', 'label', 'Гарантия', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'warranty'), 'sr', 'label', 'Garancija', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'warranty'), 'en', 'label', 'Warranty', true),

-- Период гарантии
('attribute', (SELECT id FROM category_attributes WHERE name = 'warranty_period'), 'ru', 'label', 'Срок гарантии', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'warranty_period'), 'sr', 'label', 'Period garancije', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'warranty_period'), 'en', 'label', 'Warranty period', true),

-- Возможность обмена
('attribute', (SELECT id FROM category_attributes WHERE name = 'exchange_possible'), 'ru', 'label', 'Возможен обмен', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'exchange_possible'), 'sr', 'label', 'Moguća zamena', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'exchange_possible'), 'en', 'label', 'Exchange possible', true),

-- Финансирование
('attribute', (SELECT id FROM category_attributes WHERE name = 'financing_available'), 'ru', 'label', 'Кредит/Лизинг', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'financing_available'), 'sr', 'label', 'Kredit/Lizing', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'financing_available'), 'en', 'label', 'Financing available', true),

-- Страна происхождения
('attribute', (SELECT id FROM category_attributes WHERE name = 'country_origin'), 'ru', 'label', 'Страна происхождения', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'country_origin'), 'sr', 'label', 'Zemlja porekla', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'country_origin'), 'en', 'label', 'Country of origin', true),

-- Дата первой регистрации
('attribute', (SELECT id FROM category_attributes WHERE name = 'first_registration'), 'ru', 'label', 'Первая регистрация', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'first_registration'), 'sr', 'label', 'Prva registracija', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'first_registration'), 'en', 'label', 'First registration', true),

-- Техосмотр до
('attribute', (SELECT id FROM category_attributes WHERE name = 'inspection_valid_until'), 'ru', 'label', 'Техосмотр до', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'inspection_valid_until'), 'sr', 'label', 'Tehnički pregled do', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'inspection_valid_until'), 'en', 'label', 'Inspection valid until', true),

-- Регистрация до
('attribute', (SELECT id FROM category_attributes WHERE name = 'registration_valid_until'), 'ru', 'label', 'Регистрация до', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'registration_valid_until'), 'sr', 'label', 'Registracija do', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'registration_valid_until'), 'en', 'label', 'Registration valid until', true),

-- Особенности комплектации
('attribute', (SELECT id FROM category_attributes WHERE name = 'equipment_features'), 'ru', 'label', 'Комплектация', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'equipment_features'), 'sr', 'label', 'Oprema', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'equipment_features'), 'en', 'label', 'Equipment', true),

-- Дополнительное оборудование
('attribute', (SELECT id FROM category_attributes WHERE name = 'additional_equipment'), 'ru', 'label', 'Дополнительное оборудование', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'additional_equipment'), 'sr', 'label', 'Dodatna oprema', true),
('attribute', (SELECT id FROM category_attributes WHERE name = 'additional_equipment'), 'en', 'label', 'Additional equipment', true)
) AS v(entity_type, entity_id, language, field_name, translated_text, is_verified)
WHERE v.entity_id IS NOT NULL
AND NOT EXISTS (
    SELECT 1 FROM translations t 
    WHERE t.entity_type = v.entity_type 
    AND t.entity_id = v.entity_id 
    AND t.language = v.language 
    AND t.field_name = v.field_name
);

-- Добавляем переводы для опций состояния автомобиля
INSERT INTO attribute_option_translations (attribute_name, option_value, ru_translation, sr_translation)
SELECT * FROM (VALUES
('car_condition', 'new', 'Новый', 'Nov'),
('car_condition', 'used', 'Б/у', 'Polovan'),
('car_condition', 'damaged', 'Повреждённый', 'Oštećen'),
('car_condition', 'for_parts', 'На запчасти', 'Za delove')
) AS v(attribute_name, option_value, ru_translation, sr_translation)
WHERE NOT EXISTS (
    SELECT 1 FROM attribute_option_translations aot 
    WHERE aot.attribute_name = v.attribute_name 
    AND aot.option_value = v.option_value
);

-- Добавляем переводы для опций оборудования
INSERT INTO attribute_option_translations (attribute_name, option_value, ru_translation, sr_translation)
SELECT * FROM (VALUES
('equipment_features', 'abs', 'ABS', 'ABS'),
('equipment_features', 'esp', 'ESP', 'ESP'),
('equipment_features', 'airbag', 'Подушки безопасности', 'Vazdušni jastuci'),
('equipment_features', 'climate_control', 'Климат-контроль', 'Klima uređaj'),
('equipment_features', 'cruise_control', 'Круиз-контроль', 'Tempomat'),
('equipment_features', 'parking_sensors', 'Парктроник', 'Parking senzori'),
('equipment_features', 'rear_camera', 'Камера заднего вида', 'Kamera za vožnju unazad'),
('equipment_features', 'navigation', 'Навигация', 'Navigacija'),
('equipment_features', 'leather_seats', 'Кожаный салон', 'Kožna sedišta'),
('equipment_features', 'heated_seats', 'Подогрев сидений', 'Grejanje sedišta'),
('equipment_features', 'sunroof', 'Люк', 'Šiber'),
('equipment_features', 'xenon', 'Ксенон', 'Ksenon'),
('equipment_features', 'led', 'LED фары', 'LED farovi'),
('equipment_features', 'alloy_wheels', 'Литые диски', 'Alu felne'),
('equipment_features', 'tow_bar', 'Фаркоп', 'Kuka za vuču')
) AS v(attribute_name, option_value, ru_translation, sr_translation)
WHERE NOT EXISTS (
    SELECT 1 FROM attribute_option_translations aot 
    WHERE aot.attribute_name = v.attribute_name 
    AND aot.option_value = v.option_value
);

-- Индекс для быстрого поиска по VIN
-- Примечание: индекс на text_value уже существует (idx_listing_attr_text), он будет использоваться для поиска по VIN
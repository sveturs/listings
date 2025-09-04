-- Добавление атрибутов для 7 категорий без атрибутов
-- Migration 000047: Add attributes for missing categories

-- Создание атрибутов для категории Health & Beauty (1014)
INSERT INTO unified_attributes (code, name, display_name, attribute_type, purpose, options, validation_rules, ui_settings, is_searchable, is_filterable, is_required, is_variant_compatible, affects_stock, affects_price, sort_order, is_active) VALUES
    ('health_product_type', 'Product Type', 'Тип продукта', 'select', 'regular', '["Skincare", "Makeup", "Haircare", "Perfume", "Vitamins", "Medical Device", "Personal Care", "Other"]', '{"required": true}', '{"showInCard": true}', true, true, true, false, false, false, 10, true),
    ('health_brand', 'Brand', 'Бренд', 'text', 'regular', NULL, '{"minLength": 2, "maxLength": 50}', '{"showInCard": true}', true, true, false, false, false, false, 20, true),
    ('health_condition', 'Condition', 'Состояние', 'select', 'regular', '["New", "Like New", "Good", "Fair", "Poor"]', '{"required": true}', '{"showInCard": true}', true, true, true, false, false, false, 30, true),
    ('health_expiry_date', 'Expiry Date', 'Срок годности', 'date', 'regular', NULL, NULL, '{"showInCard": false}', false, false, false, false, false, false, 40, true),
    ('health_skin_type', 'Skin Type', 'Тип кожи', 'multiselect', 'regular', '["Dry", "Oily", "Combination", "Sensitive", "Normal", "All Types"]', NULL, '{"showInCard": false}', false, true, false, true, false, false, 50, true);

-- Связывание атрибутов с категорией Health & Beauty
INSERT INTO unified_category_attributes (category_id, attribute_id, is_enabled, is_required, sort_order) 
SELECT 1014, id, true, (code IN ('health_product_type', 'health_condition')), sort_order
FROM unified_attributes WHERE code IN ('health_product_type', 'health_brand', 'health_condition', 'health_expiry_date', 'health_skin_type');

-- Создание атрибутов для категории Kids & Baby (1013)
INSERT INTO unified_attributes (code, name, display_name, attribute_type, purpose, options, validation_rules, ui_settings, is_searchable, is_filterable, is_required, is_variant_compatible, affects_stock, affects_price, sort_order, is_active) VALUES
    ('kids_age_group', 'Age Group', 'Возрастная группа', 'select', 'regular', '["0-6 months", "6-12 months", "1-2 years", "3-5 years", "6-8 years", "9-12 years", "13+ years"]', '{"required": true}', '{"showInCard": true}', true, true, true, false, false, false, 10, true),
    ('kids_product_type', 'Product Type', 'Тип продукта', 'select', 'regular', '["Toys", "Clothing", "Shoes", "Baby Gear", "Books", "Furniture", "Safety", "Educational", "Other"]', '{"required": true}', '{"showInCard": true}', true, true, true, false, false, false, 20, true),
    ('kids_brand', 'Brand', 'Бренд', 'text', 'regular', NULL, '{"minLength": 2, "maxLength": 50}', '{"showInCard": true}', true, true, false, false, false, false, 30, true),
    ('kids_condition', 'Condition', 'Состояние', 'select', 'regular', '["New", "Like New", "Good", "Fair", "Poor"]', '{"required": true}', '{"showInCard": true}', true, true, true, false, false, false, 40, true),
    ('kids_size', 'Size', 'Размер', 'text', 'regular', NULL, '{"maxLength": 20}', '{"showInCard": false}', true, true, false, true, false, false, 50, true),
    ('kids_gender', 'Gender', 'Пол', 'select', 'regular', '["Boy", "Girl", "Unisex"]', NULL, '{"showInCard": false}', false, true, false, false, false, false, 60, true);

INSERT INTO unified_category_attributes (category_id, attribute_id, is_enabled, is_required, sort_order) 
SELECT 1013, id, true, (code IN ('kids_age_group', 'kids_product_type', 'kids_condition')), sort_order
FROM unified_attributes WHERE code IN ('kids_age_group', 'kids_product_type', 'kids_brand', 'kids_condition', 'kids_size', 'kids_gender');

-- Создание атрибутов для категории Musical Instruments (1016)
INSERT INTO unified_attributes (code, name, display_name, attribute_type, purpose, options, validation_rules, ui_settings, is_searchable, is_filterable, is_required, is_variant_compatible, affects_stock, affects_price, sort_order, is_active) VALUES
    ('music_instrument_type', 'Instrument Type', 'Тип инструмента', 'select', 'regular', '["Guitar", "Piano", "Drums", "Violin", "Bass", "Keyboard", "Microphone", "Amplifier", "Effects", "DJ Equipment", "Other"]', '{"required": true}', '{"showInCard": true}', true, true, true, false, false, false, 10, true),
    ('music_brand', 'Brand', 'Бренд', 'text', 'regular', NULL, '{"minLength": 2, "maxLength": 50}', '{"showInCard": true}', true, true, false, false, false, false, 20, true),
    ('music_condition', 'Condition', 'Состояние', 'select', 'regular', '["New", "Like New", "Good", "Fair", "Poor"]', '{"required": true}', '{"showInCard": true}', true, true, true, false, false, false, 30, true),
    ('music_skill_level', 'Skill Level', 'Уровень игры', 'select', 'regular', '["Beginner", "Intermediate", "Professional", "All Levels"]', NULL, '{"showInCard": false}', false, true, false, false, false, false, 40, true),
    ('music_acoustic_electric', 'Type', 'Тип', 'select', 'regular', '["Acoustic", "Electric", "Hybrid", "Digital", "Not Applicable"]', NULL, '{"showInCard": false}', true, true, false, false, false, false, 50, true);

INSERT INTO unified_category_attributes (category_id, attribute_id, is_enabled, is_required, is_filter, sort_order, show_in_card, show_in_list) 
SELECT 1016, id, true, (code IN ('music_instrument_type', 'music_condition')), true, sort_order, (code IN ('music_instrument_type', 'music_brand', 'music_condition')), false
FROM unified_attributes WHERE code IN ('music_instrument_type', 'music_brand', 'music_condition', 'music_skill_level', 'music_acoustic_electric');

-- Создание атрибутов для категории Events & Tickets (1020)
INSERT INTO unified_attributes (code, name, display_name, attribute_type, purpose, options, validation_rules, ui_settings, is_searchable, is_filterable, is_required, is_variant_compatible, affects_stock, affects_price, sort_order, is_active) VALUES
    ('event_type', 'Event Type', 'Тип мероприятия', 'select', 'regular', '["Concert", "Theater", "Sports", "Conference", "Festival", "Exhibition", "Workshop", "Party", "Travel", "Other"]', '{"required": true}', '{"showInCard": true}', true, true, true, false, false, false, 10, true),
    ('event_date', 'Event Date', 'Дата мероприятия', 'date', 'regular', NULL, '{"required": true}', '{"showInCard": true}', true, true, true, false, false, false, 20, true),
    ('event_location', 'Location', 'Место проведения', 'text', 'regular', NULL, '{"minLength": 3, "maxLength": 100, "required": true}', '{"showInCard": true}', true, true, true, false, false, false, 30, true),
    ('event_quantity', 'Ticket Quantity', 'Количество билетов', 'number', 'regular', NULL, '{"min": 1, "max": 100, "required": true}', '{"showInCard": true}', false, false, true, false, true, false, 40, true),
    ('event_seat_type', 'Seat Type', 'Тип мест', 'select', 'regular', '["VIP", "Premium", "Standard", "General Admission", "Standing", "Other"]', NULL, '{"showInCard": false}', false, true, false, false, false, false, 50, true);

INSERT INTO unified_category_attributes (category_id, attribute_id, is_enabled, is_required, is_filter, sort_order, show_in_card, show_in_list) 
SELECT 1020, id, true, (code IN ('event_type', 'event_date', 'event_location', 'event_quantity')), true, sort_order, (code IN ('event_type', 'event_date', 'event_location', 'event_quantity')), false
FROM unified_attributes WHERE code IN ('event_type', 'event_date', 'event_location', 'event_quantity', 'event_seat_type');

-- Создание атрибутов для категории Education (1019)
INSERT INTO unified_attributes (code, name, display_name, attribute_type, purpose, options, validation_rules, ui_settings, is_searchable, is_filterable, is_required, is_variant_compatible, affects_stock, affects_price, sort_order, is_active) VALUES
    ('edu_type', 'Education Type', 'Тип образования', 'select', 'regular', '["Course", "Tutorial", "Book", "Software", "Certification", "Workshop", "Mentoring", "Equipment", "Other"]', '{"required": true}', '{"showInCard": true}', true, true, true, false, false, false, 10, true),
    ('edu_subject', 'Subject', 'Предмет', 'text', 'regular', NULL, '{"minLength": 2, "maxLength": 50, "required": true}', '{"showInCard": true}', true, true, true, false, false, false, 20, true),
    ('edu_level', 'Level', 'Уровень', 'select', 'regular', '["Beginner", "Intermediate", "Advanced", "Expert", "All Levels"]', '{"required": true}', '{"showInCard": true}', true, true, true, false, false, false, 30, true),
    ('edu_language', 'Language', 'Язык', 'select', 'regular', '["Serbian", "English", "Russian", "German", "French", "Spanish", "Italian", "Other"]', '{"required": true}', '{"showInCard": false}', false, true, true, false, false, false, 40, true),
    ('edu_format', 'Format', 'Формат', 'select', 'regular', '["Online", "In-Person", "Hybrid", "Self-Paced", "Live Session"]', NULL, '{"showInCard": false}', false, true, false, false, false, false, 50, true);

INSERT INTO unified_category_attributes (category_id, attribute_id, is_enabled, is_required, is_filter, sort_order, show_in_card, show_in_list) 
SELECT 1019, id, true, (code IN ('edu_type', 'edu_subject', 'edu_level', 'edu_language')), true, sort_order, (code IN ('edu_type', 'edu_subject', 'edu_level')), false
FROM unified_attributes WHERE code IN ('edu_type', 'edu_subject', 'edu_level', 'edu_language', 'edu_format');

-- Создание атрибутов для категории Antiques & Art (1017)
INSERT INTO unified_attributes (code, name, display_name, attribute_type, purpose, options, validation_rules, ui_settings, is_searchable, is_filterable, is_required, is_variant_compatible, affects_stock, affects_price, sort_order, is_active) VALUES
    ('art_type', 'Art Type', 'Тип искусства', 'select', 'regular', '["Painting", "Sculpture", "Photography", "Print", "Drawing", "Antique", "Collectible", "Jewelry", "Furniture", "Ceramics", "Other"]', '{"required": true}', '{"showInCard": true}', true, true, true, false, false, false, 10, true),
    ('art_period', 'Period/Era', 'Период/Эпоха', 'text', 'regular', NULL, '{"minLength": 3, "maxLength": 50}', '{"showInCard": true}', true, true, false, false, false, false, 20, true),
    ('art_condition', 'Condition', 'Состояние', 'select', 'regular', '["Excellent", "Very Good", "Good", "Fair", "Restoration Needed"]', '{"required": true}', '{"showInCard": true}', true, true, true, false, false, false, 30, true),
    ('art_authenticity', 'Authenticity', 'Подлинность', 'select', 'regular', '["Certified Original", "Original", "Reproduction", "Print", "Unknown"]', '{"required": true}', '{"showInCard": true}', true, true, true, false, false, false, 40, true),
    ('art_material', 'Material', 'Материал', 'text', 'regular', NULL, '{"minLength": 2, "maxLength": 100}', '{"showInCard": false}', true, true, false, false, false, false, 50, true),
    ('art_dimensions', 'Dimensions', 'Размеры', 'text', 'regular', NULL, '{"maxLength": 50}', '{"showInCard": false}', true, false, false, false, false, false, 60, true);

INSERT INTO unified_category_attributes (category_id, attribute_id, is_enabled, is_required, is_filter, sort_order, show_in_card, show_in_list) 
SELECT 1017, id, true, (code IN ('art_type', 'art_condition', 'art_authenticity')), true, sort_order, (code IN ('art_type', 'art_period', 'art_condition', 'art_authenticity')), false
FROM unified_attributes WHERE code IN ('art_type', 'art_period', 'art_condition', 'art_authenticity', 'art_material', 'art_dimensions');

-- Создание атрибутов для категории Hobbies & Entertainment (1015)
INSERT INTO unified_attributes (code, name, display_name, attribute_type, purpose, options, validation_rules, ui_settings, is_searchable, is_filterable, is_required, is_variant_compatible, affects_stock, affects_price, sort_order, is_active) VALUES
    ('hobby_type', 'Hobby Type', 'Тип хобби', 'select', 'regular', '["Board Games", "Video Games", "Puzzles", "Models", "Crafts", "Sports Equipment", "Outdoor Gear", "Books", "Movies", "Collectibles", "Other"]', '{"required": true}', '{"showInCard": true}', true, true, true, false, false, false, 10, true),
    ('hobby_brand', 'Brand', 'Бренд', 'text', 'regular', NULL, '{"minLength": 2, "maxLength": 50}', '{"showInCard": true}', true, true, false, false, false, false, 20, true),
    ('hobby_condition', 'Condition', 'Состояние', 'select', 'regular', '["New", "Like New", "Good", "Fair", "Poor"]', '{"required": true}', '{"showInCard": true}', true, true, true, false, false, false, 30, true),
    ('hobby_age_group', 'Age Group', 'Возрастная группа', 'select', 'regular', '["Kids (3-12)", "Teens (13-17)", "Adults (18+)", "All Ages"]', NULL, '{"showInCard": false}', false, true, false, false, false, false, 40, true),
    ('hobby_skill_level', 'Skill Level', 'Уровень сложности', 'select', 'regular', '["Beginner", "Intermediate", "Advanced", "Expert", "All Levels"]', NULL, '{"showInCard": false}', false, true, false, false, false, false, 50, true),
    ('hobby_players', 'Number of Players', 'Количество игроков', 'text', 'regular', NULL, '{"maxLength": 20}', '{"showInCard": false}', false, false, false, false, false, false, 60, true);

INSERT INTO unified_category_attributes (category_id, attribute_id, is_enabled, is_required, is_filter, sort_order, show_in_card, show_in_list) 
SELECT 1015, id, true, (code IN ('hobby_type', 'hobby_condition')), true, sort_order, (code IN ('hobby_type', 'hobby_brand', 'hobby_condition')), false
FROM unified_attributes WHERE code IN ('hobby_type', 'hobby_brand', 'hobby_condition', 'hobby_age_group', 'hobby_skill_level', 'hobby_players');

-- Обновление статистики атрибутов
INSERT INTO unified_attribute_stats (category_id, total_attributes, active_attributes, last_updated)
VALUES 
    (1014, 5, 5, NOW()), -- Health & Beauty
    (1013, 6, 6, NOW()), -- Kids & Baby  
    (1016, 5, 5, NOW()), -- Musical Instruments
    (1020, 5, 5, NOW()), -- Events & Tickets
    (1019, 5, 5, NOW()), -- Education
    (1017, 6, 6, NOW()), -- Antiques & Art
    (1015, 6, 6, NOW())  -- Hobbies & Entertainment
ON CONFLICT (category_id) DO UPDATE SET
    total_attributes = EXCLUDED.total_attributes,
    active_attributes = EXCLUDED.active_attributes,
    last_updated = NOW();
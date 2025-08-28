-- Добавление новых основных категорий
INSERT INTO marketplace_categories (id, slug, parent_id, icon, created_at, name) VALUES
(1011, 'pets', NULL, '🐾', CURRENT_TIMESTAMP, 'Pets'),
(1012, 'books-stationery', NULL, '📚', CURRENT_TIMESTAMP, 'Books & Stationery'),
(1013, 'kids-baby', NULL, '👶', CURRENT_TIMESTAMP, 'Kids & Baby'),
(1014, 'health-beauty', NULL, '💄', CURRENT_TIMESTAMP, 'Health & Beauty'),
(1015, 'hobbies-entertainment', NULL, '🎮', CURRENT_TIMESTAMP, 'Hobbies & Entertainment'),
(1016, 'musical-instruments', NULL, '🎸', CURRENT_TIMESTAMP, 'Musical Instruments'),
(1017, 'antiques-art', NULL, '🎨', CURRENT_TIMESTAMP, 'Antiques & Art'),
(1018, 'jobs', NULL, '💼', CURRENT_TIMESTAMP, 'Jobs'),
(1019, 'education', NULL, '🎓', CURRENT_TIMESTAMP, 'Education'),
(1020, 'events-tickets', NULL, '🎫', CURRENT_TIMESTAMP, 'Events & Tickets');

-- Добавление подкатегорий для электроники
INSERT INTO marketplace_categories (id, slug, parent_id, icon, created_at, name) VALUES
(1105, 'gaming-consoles', 1001, '🎮', CURRENT_TIMESTAMP, 'Gaming Consoles'),
(1106, 'photo-video', 1001, '📷', CURRENT_TIMESTAMP, 'Photo & Video'),
(1107, 'smart-home', 1001, '🏠', CURRENT_TIMESTAMP, 'Smart Home'),
(1108, 'electronics-accessories', 1001, '🔌', CURRENT_TIMESTAMP, 'Electronics Accessories');

-- Добавление подкатегорий для моды
INSERT INTO marketplace_categories (id, slug, parent_id, icon, created_at, name) VALUES
(1205, 'kids-clothing', 1002, '👕', CURRENT_TIMESTAMP, 'Kids Clothing'),
(1206, 'sports-clothing', 1002, '🏃', CURRENT_TIMESTAMP, 'Sports Clothing'),
(1207, 'watches', 1002, '⌚', CURRENT_TIMESTAMP, 'Watches'),
(1208, 'bags', 1002, '👜', CURRENT_TIMESTAMP, 'Bags');

-- Добавление подкатегорий для дома и сада
INSERT INTO marketplace_categories (id, slug, parent_id, icon, created_at, name) VALUES
(1505, 'kitchenware', 1005, '🍽️', CURRENT_TIMESTAMP, 'Kitchenware'),
(1506, 'textiles', 1005, '🛏️', CURRENT_TIMESTAMP, 'Textiles'),
(1507, 'lighting', 1005, '💡', CURRENT_TIMESTAMP, 'Lighting'),
(1508, 'plumbing', 1005, '🚿', CURRENT_TIMESTAMP, 'Plumbing');

-- Удаляем тестовую категорию
DELETE FROM marketplace_categories WHERE id = 2005;

-- Добавление подкатегорий для спорта
INSERT INTO marketplace_categories (id, slug, parent_id, icon, created_at, name) VALUES
(2011, 'bicycles', 1010, '🚴', CURRENT_TIMESTAMP, 'Bicycles'),
(2012, 'water-sports', 1010, '🏊', CURRENT_TIMESTAMP, 'Water Sports'),
(2013, 'hunting-fishing', 1010, '🎣', CURRENT_TIMESTAMP, 'Hunting & Fishing'),
(2014, 'camping-hiking', 1010, '🏕️', CURRENT_TIMESTAMP, 'Camping & Hiking');

-- Добавление новых общих атрибутов
INSERT INTO category_attributes (id, name, display_name, attribute_type, options, is_searchable, is_filterable, is_required, sort_order) VALUES
(2601, 'location', 'Lokacija', 'text', '{}', true, true, false, 1),
(2602, 'delivery_available', 'Dostava', 'boolean', '{}', true, true, false, 50),
(2603, 'negotiable', 'Dogovor', 'boolean', '{}', true, true, false, 51),
(2604, 'warranty', 'Garancija', 'select', '{"values": ["no_warranty", "manufacturer", "store", "extended"]}', true, true, false, 52),
(2605, 'return_policy', 'Povrat', 'select', '{"values": ["no_returns", "7_days", "14_days", "30_days"]}', true, true, false, 53);

-- Атрибуты для электроники
INSERT INTO category_attributes (id, name, display_name, attribute_type, options, is_searchable, is_filterable, is_required, sort_order) VALUES
(2701, 'screen_size', 'Veličina ekrana', 'select', '{"values": ["5\"", "6\"", "7\"", "8\"", "10\"", "11\"", "13\"", "15\"", "17\"", "21\"", "24\"", "27\"", "32\"", "40\"", "43\"", "50\"", "55\"", "65\"", "75\""]}', true, true, false, 10),
(2702, 'battery_life', 'Trajanje baterije', 'select', '{"values": ["<4h", "4-8h", "8-12h", "12-24h", "24-48h", ">48h"]}', true, true, false, 11),
(2703, 'connectivity', 'Povezivanje', 'multiselect', '{"values": ["wifi", "bluetooth", "nfc", "usb_c", "hdmi", "ethernet", "3g", "4g", "5g"]}', true, true, false, 12),
(2704, 'resolution', 'Rezolucija', 'select', '{"values": ["HD", "Full HD", "2K", "4K", "8K", "720p", "1080p", "1440p", "2160p"]}', true, true, false, 13);

-- Атрибуты для моды
INSERT INTO category_attributes (id, name, display_name, attribute_type, options, is_searchable, is_filterable, is_required, sort_order) VALUES
(2801, 'size', 'Veličina', 'select', '{"values": ["XS", "S", "M", "L", "XL", "XXL", "XXXL", "36", "37", "38", "39", "40", "41", "42", "43", "44", "45", "46"]}', true, true, false, 10),
(2802, 'material', 'Materijal', 'multiselect', '{"values": ["cotton", "polyester", "wool", "silk", "leather", "synthetic", "denim", "linen"]}', true, true, false, 11),
(2803, 'gender', 'Pol', 'select', '{"values": ["male", "female", "unisex", "boys", "girls"]}', true, true, false, 12),
(2804, 'season', 'Sezona', 'multiselect', '{"values": ["spring", "summer", "autumn", "winter", "all_seasons"]}', true, true, false, 13);

-- Атрибуты для недвижимости
INSERT INTO category_attributes (id, name, display_name, attribute_type, options, is_searchable, is_filterable, is_required, sort_order) VALUES
(2901, 'heating_type', 'Tip grejanja', 'select', '{"values": ["central", "electric", "gas", "solid_fuel", "heat_pump", "floor_heating", "no_heating"]}', true, true, false, 20),
(2902, 'construction_year', 'Godina izgradnje', 'number', '{"min": 1900, "max": 2024}', true, true, false, 21),
(2903, 'elevator', 'Lift', 'boolean', '{}', true, true, false, 22),
(2904, 'security', 'Obezbeđenje', 'boolean', '{}', true, true, false, 23);

-- Атрибуты для автомобилей
INSERT INTO category_attributes (id, name, display_name, attribute_type, options, is_searchable, is_filterable, is_required, sort_order) VALUES
(3001, 'engine_size', 'Zapremina motora', 'select', '{"values": ["<1.0L", "1.0L", "1.2L", "1.4L", "1.6L", "1.8L", "2.0L", "2.5L", "3.0L", ">3.0L"]}', true, true, false, 20),
(3002, 'doors', 'Broj vrata', 'select', '{"values": ["2", "3", "4", "5"]}', true, true, false, 21),
(3003, 'seats', 'Broj sedišta', 'select', '{"values": ["2", "4", "5", "6", "7", "8", "9+"]}', true, true, false, 22),
(3004, 'drive_type', 'Pogon', 'select', '{"values": ["front", "rear", "awd", "4wd"]}', true, true, false, 23);

-- Атрибуты для услуг
INSERT INTO category_attributes (id, name, display_name, attribute_type, options, is_searchable, is_filterable, is_required, sort_order) VALUES
(3101, 'price_type', 'Tip cene', 'select', '{"values": ["per_hour", "per_day", "per_project", "per_month", "fixed", "negotiable"]}', true, true, false, 10),
(3102, 'certification', 'Sertifikati', 'text', '{}', true, false, false, 11),
(3103, 'language', 'Jezici', 'multiselect', '{"values": ["serbian", "english", "german", "russian", "french", "italian", "spanish", "chinese"]}', true, true, false, 12),
(3104, 'portfolio_url', 'Portfolio', 'text', '{}', false, false, false, 13);

-- Обновление пустых опций для существующих атрибутов
UPDATE category_attributes SET options = '{"values": ["sedan", "hatchback", "suv", "coupe", "wagon", "minivan", "pickup", "convertible"]}' WHERE name = 'body_type';
UPDATE category_attributes SET options = '{"values": ["consulting", "repair", "installation", "maintenance", "cleaning", "transport", "design", "development", "education", "other"]}' WHERE name = 'service_type';
UPDATE category_attributes SET options = '{"values": ["immediate", "within_24h", "within_week", "by_appointment", "weekdays", "weekends", "24_7"]}' WHERE name = 'availability';
UPDATE category_attributes SET options = '{"values": ["local", "city", "region", "country", "international", "online"]}' WHERE name = 'service_area';

-- Связывание общих атрибутов со всеми категориями
DO $$
DECLARE
    cat_id INTEGER;
BEGIN
    FOR cat_id IN SELECT id FROM marketplace_categories WHERE parent_id IS NOT NULL
    LOOP
        -- Добавляем общие атрибуты ко всем подкатегориям
        INSERT INTO category_attribute_mapping (category_id, attribute_id, is_required, sort_order)
        VALUES 
            (cat_id, 2601, false, 100), -- location
            (cat_id, 2602, false, 101), -- delivery_available
            (cat_id, 2603, false, 102), -- negotiable
            (cat_id, 2604, false, 103), -- warranty
            (cat_id, 2605, false, 104)  -- return_policy
        ON CONFLICT DO NOTHING;
    END LOOP;
END $$;

-- Связывание специфичных атрибутов с категориями электроники
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_required, sort_order)
SELECT c.id, a.id, false, 20 + (a.id - 2700)
FROM marketplace_categories c
CROSS JOIN category_attributes a
WHERE c.parent_id = 1001 AND a.id BETWEEN 2701 AND 2704
ON CONFLICT DO NOTHING;

-- Связывание специфичных атрибутов с категориями моды
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_required, sort_order)
SELECT c.id, a.id, false, 20 + (a.id - 2800)
FROM marketplace_categories c
CROSS JOIN category_attributes a
WHERE c.parent_id = 1002 AND a.id BETWEEN 2801 AND 2804
ON CONFLICT DO NOTHING;

-- Связывание специфичных атрибутов с категориями недвижимости
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_required, sort_order)
SELECT c.id, a.id, false, 30 + (a.id - 2900)
FROM marketplace_categories c
CROSS JOIN category_attributes a
WHERE c.parent_id = 1004 AND a.id BETWEEN 2901 AND 2904
ON CONFLICT DO NOTHING;

-- Связывание специфичных атрибутов с категориями автомобилей
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_required, sort_order)
SELECT c.id, a.id, false, 30 + (a.id - 3000)
FROM marketplace_categories c
CROSS JOIN category_attributes a
WHERE c.parent_id = 1003 AND a.id BETWEEN 3001 AND 3004
ON CONFLICT DO NOTHING;

-- Связывание специфичных атрибутов с категориями услуг
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_required, sort_order)
SELECT c.id, a.id, false, 20 + (a.id - 3100)
FROM marketplace_categories c
CROSS JOIN category_attributes a
WHERE c.parent_id = 1009 AND a.id BETWEEN 3101 AND 3104
ON CONFLICT DO NOTHING;